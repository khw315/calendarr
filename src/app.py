#!/usr/bin/env python3
# src/app.py

import sys
import io
from contextlib import redirect_stdout, redirect_stderr
import traceback
import datetime
import os
import logging

from flask import Flask, send_from_directory
from flask import cli
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.cron import CronTrigger
import main

from constants import (
    DEFAULT_LOG_DIR, JOB_ID_DEBUG_PING, JOB_ID_LOG_CLEANUP, JOB_ID_MAIN
)

os.makedirs(DEFAULT_LOG_DIR, exist_ok=True)


from config.settings import load_config_from_env
from utils.logging_utils import setup_logging, cleanup_log_files

app = Flask(__name__)

# Disable the annoying Flask "dev server only" banner
cli = sys.modules['flask.cli']
cli.show_server_banner = lambda *_: None  

config = None
logger = None
scheduler = None


def init_app():
    """Initialize application configuration and logging"""
    global config, logger
    
    config = load_config_from_env()

    logger = setup_logging(
        log_dir=config.logging_settings.log_dir,
        log_file=config.logging_settings.log_file,
        debug_mode=config.logging_settings.debug_mode,
        max_size_mb=config.logging_settings.max_size_mb,
        backup_count=config.logging_settings.backup_count
    )
    
    if config.logging_settings.debug_mode:
        logger.setLevel("DEBUG")
        logger.info("🪲 Running in DEBUG mode")
        
    logger.info("🚀  Application initialized with configuration")
    return config


def run_main_job():
    """Run the main calandarr script"""
    logger.info("⚙️  Running main job")

    try:
        # Capture output from main function
        output = io.StringIO()
        with redirect_stdout(output), redirect_stderr(output):
            success = main.main()

        if not success:
            logger.info(f"⚠️  Job completed with success: {success}")
        else:
            logger.info(f"✅  Job completed with success: {success}")
        
        
    except Exception as e:
        logger.error(f"☠️  Error in main job: {str(e)}")
        logger.error(traceback.format_exc())
    
    logger.info("✅  Job function complete, container should stay running")

def log_ping():
    """Log a ping message every minute when debug is enabled"""
    if os.environ.get("DEBUG", "").lower() == "true":
        logger.debug("🔄  Ping - Application is running")


# Flask routes
@app.route('/')
def index():
    """Serve the web UI"""
    public_dir = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'public')
    return send_from_directory(public_dir, 'index.html')

@app.route('/<path:path>')
def serve_static(path):
    """Serve static files"""
    public_dir = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'public')
    return send_from_directory(public_dir, path)

@app.route('/health')
def health():
    """Health check endpoint for Docker"""
    return {"status": "healthy", "timestamp": datetime.datetime.now().isoformat()}

@app.route('/api/releases')
def get_events():
    """Get upcoming releases"""
    from flask import jsonify, request
    try:
        # Get days parameter from query string (default to 2 for today + tomorrow)
        days = request.args.get('days', default=2, type=int)
        
        # Limit to reasonable range (1-30 days)
        days = max(1, min(days, 30))
        
        # For Web UI, use customizable date range
        from datetime import datetime, timedelta
        import pytz
        
        tz = config.timezone_obj
        now = datetime.now(tz)
        start_date = now.replace(hour=0, minute=0, second=0, microsecond=0)
        end_date = start_date + timedelta(days=days)
        
        logger.info(f"Web UI requesting events from {start_date.date()} to {end_date.date()} ({days} days)")
        
        # Get events
        from services.calendar_service import CalendarService
        from services.formatter_service import FormatterService
        
        calendar_service = CalendarService(config)
        formatter_service = FormatterService(config)
        
        events = calendar_service.fetch_events(start_date, end_date)
        days, events_summary = formatter_service.process_events(
            events, 
            start_date, 
            end_date
        )
        
        # Format response
        days_data = []
        for day in days:
            day_events = []
            
            # Group TV events by timestamp and show_name
            tv_groups = {}
            for event_item in day.tv_events:
                if event_item.is_past:
                    continue
                
                show_name = event_item.show_name if getattr(event_item, 'show_name', None) else (event_item.title if hasattr(event_item, 'title') else event_item.summary)
                key = (getattr(event_item, 'timestamp', None), show_name)
                
                if key not in tv_groups:
                    tv_groups[key] = []
                tv_groups[key].append(event_item)

            BULK_THRESHOLD = 2

            for key, group in tv_groups.items():
                if len(group) >= BULK_THRESHOLD:
                    # Bulk format
                    first = group[0]
                    show_name = key[1]
                    is_premiere = any(getattr(e, 'is_premiere', False) for e in group)

                    title = f"{show_name}"
                    if is_premiere:
                        title += " 🎉"
                        
                    day_events.append({
                        'title': title,
                        'start_time': first.time_str if hasattr(first, 'time_str') else None,
                        'timestamp': first.timestamp if hasattr(first, 'timestamp') else None,
                        'end_time': first.end_time_str if hasattr(first, 'end_time_str') else None,
                        'end_timestamp': first.end_timestamp if hasattr(first, 'end_timestamp') else None,
                        'type': 'tv',
                        'is_bulk': True,
                        'episode_count': len(group)
                    })
                else:
                    # Individual format
                    for event_item in group:
                        day_events.append({
                            'title': event_item.title if hasattr(event_item, 'title') else event_item.summary,
                            'start_time': event_item.time_str if hasattr(event_item, 'time_str') else None,
                            'timestamp': event_item.timestamp if hasattr(event_item, 'timestamp') else None,
                            'end_time': event_item.end_time_str if hasattr(event_item, 'end_time_str') else None,
                            'end_timestamp': event_item.end_timestamp if hasattr(event_item, 'end_timestamp') else None,
                            'type': 'tv'
                        })
            
            for event_item in day.movie_events:
                # Always skip past events in the Web UI
                if event_item.is_past:
                    continue
                    
                day_events.append({
                    'title': event_item.title if hasattr(event_item, 'title') else event_item.summary,
                    'start_time': event_item.time_str if hasattr(event_item, 'time_str') else None,
                    'timestamp': event_item.timestamp if hasattr(event_item, 'timestamp') else None,
                    'end_time': event_item.end_time_str if hasattr(event_item, 'end_time_str') else None,
                    'end_timestamp': event_item.end_timestamp if hasattr(event_item, 'end_timestamp') else None,
                    'type': 'movie'
                })
            
            # Only include days that have events (after filtering)
            if day_events:
                days_data.append({
                    'day_name': day.day_name if day.day_name else day.name.split(',')[0],
                    'date': day.date.strftime('%B %d, %Y') if day.date else day.name,
                    'events': day_events
                })
        
        return jsonify({
            'days': days_data,
            'total_events': len(events),
            'start_date': start_date.strftime('%Y-%m-%d'),
            'end_date': end_date.strftime('%Y-%m-%d')
        })
    except Exception as e:
        logger.error(f"Error getting events: {e}")
        import traceback
        logger.error(traceback.format_exc())
        return jsonify({'error': str(e)}), 500



@app.route('/api/past-releases')
def get_past_events():
    """Get past releases"""
    from flask import jsonify, request
    try:
        # Get last 7 days of events
        from datetime import datetime, timedelta
        import pytz
        
        
        # Get days parameter from query string (default to 7)
        days = request.args.get('days', default=7, type=int)
        days = max(1, min(days, 30))
        
        tz = config.timezone_obj
        now = datetime.now(tz)
        start_date = now.replace(hour=0, minute=0, second=0, microsecond=0) - timedelta(days=days)
        end_date = now  # Up to current time
        
        logger.info(f"Web UI requesting past events from {start_date} to {end_date}")
        
        # Get events
        from services.calendar_service import CalendarService
        from services.formatter_service import FormatterService
        
        calendar_service = CalendarService(config)
        formatter_service = FormatterService(config)
        
        events = calendar_service.fetch_events(start_date, end_date)
        days, events_summary = formatter_service.process_events(
            events, 
            start_date, 
            end_date
        )
        
        # Collect only past events
        past_events = []
        for day in days:
            # Group TV events by timestamp and show_name
            tv_groups = {}
            for event_item in day.tv_events:
                if event_item.is_past:
                    show_name = event_item.show_name if getattr(event_item, 'show_name', None) else (event_item.title if hasattr(event_item, 'title') else event_item.summary)
                    key = (getattr(event_item, 'timestamp', None), show_name)
                    
                    if key not in tv_groups:
                        tv_groups[key] = []
                    tv_groups[key].append(event_item)

            BULK_THRESHOLD = 4

            for key, group in tv_groups.items():
                if len(group) >= BULK_THRESHOLD:
                    # Bulk format
                    first = group[0]
                    show_name = key[1]
                    is_premiere = any(getattr(e, 'is_premiere', False) for e in group)
                    
                    title = f"{show_name}"
                    if is_premiere:
                        title += " 🎉"
                        
                    past_events.append({
                        'title': title,
                        'start_time': first.time_str if hasattr(first, 'time_str') else None,
                        'date': day.date.strftime('%B %d, %Y') if day.date else day.name,
                        'type': 'tv',
                        'timestamp': getattr(first, 'timestamp', None),
                        'is_bulk': True,
                        'episode_count': len(group)
                    })
                else:
                    # Individual format
                    for event_item in group:
                        past_events.append({
                            'title': event_item.title if hasattr(event_item, 'title') else event_item.summary,
                            'start_time': event_item.time_str if hasattr(event_item, 'time_str') else None,
                            'date': day.date.strftime('%B %d, %Y') if day.date else day.name,
                            'type': 'tv',
                            'timestamp': getattr(event_item, 'timestamp', None)
                        })
            
            # Check movie events
            for event_item in day.movie_events:
                if event_item.is_past:
                    past_events.append({
                        'title': event_item.title if hasattr(event_item, 'title') else event_item.summary,
                        'start_time': event_item.time_str if hasattr(event_item, 'time_str') else None,
                        'date': day.date.strftime('%B %d, %Y') if day.date else day.name,
                        'type': 'movie',
                        'timestamp': getattr(event_item, 'timestamp', None)
                    })
        
        return jsonify({
            'events': past_events,
            'count': len(past_events)
        })
    except Exception as e:
        logger.error(f"Error getting past events: {e}")
        import traceback
        logger.error(traceback.format_exc())
        return jsonify({'error': str(e)}), 500

@app.route('/api/schedule')
def get_schedule():
    """Get schedule information"""
    from flask import jsonify
    try:
        schedule = config.schedule_settings
        
        # Calculate next run time
        next_run = None
        if scheduler:
            job = scheduler.get_job(JOB_ID_MAIN)
            if job and job.next_run_time:
                next_run = job.next_run_time.isoformat()
        
        return jsonify({
            'schedule_type': schedule.schedule_type,
            'run_time': schedule.run_time,
            'next_run': next_run,
            'cron_schedule': schedule.cron_schedule,
            'timezone': config.timezone
        })
    except Exception as e:
        logger.error(f"Error getting schedule: {e}")
        return jsonify({'error': str(e)}), 500

@app.route('/api/trigger', methods=['POST'])
def trigger_job():
    """Manually trigger the calendar job"""
    from flask import jsonify
    try:
        logger.info("Manual trigger requested via API")
        run_main_job()
        return jsonify({
            'success': True,
            'message': 'Job triggered successfully'
        })
    except Exception as e:
        logger.error(f"Error triggering job: {e}")
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500

@app.route('/api/config', methods=['GET', 'POST'])
def handle_config():
    """Get or save configuration"""
    global config, scheduler
    from flask import jsonify, request
    from config.settings import load_config_from_file, save_config_to_file
    
    if request.method == 'GET':
        try:
            saved_config = load_config_from_file()

            current_config = {
                # General
                'DEBUG': config.logging_settings.debug_mode,
                
                # Platforms
                'USE_DISCORD': config.use_discord,
                'DISCORD_WEBHOOK_URL': config.discord_webhook_url or "",
                'DISCORD_MENTION_ROLE_ID': config.discord_mention_role_id or "",
                'DISCORD_HIDE_MENTION_INSTRUCTIONS': config.discord_hide_mention_instructions,
                'DISCORD_TIMESTAMP_STYLE': config.discord_timestamp_style or "Relative Time",
                'ENABLE_CUSTOM_DISCORD_FOOTER': config.enable_custom_discord_footer,
                
                'USE_SLACK': config.use_slack,
                'SLACK_WEBHOOK_URL': config.slack_webhook_url or "",
                'ENABLE_CUSTOM_SLACK_FOOTER': config.enable_custom_slack_footer,
                
                # Calendar
                'CALENDAR_URLS': [url.to_dict() for url in config.calendar_urls],
                'PASSED_EVENT_HANDLING': config.passed_event_handling,
                'DEDUPLICATE_EVENTS': config.deduplicate_events,
                
                # Time
                'USE_24_HOUR': config.time_settings.use_24_hour,
                'ADD_LEADING_ZERO': config.time_settings.add_leading_zero,
                'DISPLAY_TIME': config.time_settings.display_time,
                'SHOW_DATE_RANGE': config.show_date_range,
                'SHOW_TIMEZONE_IN_SUBHEADER': config.show_timezone_in_subheader,
                'TZ': config.timezone,
                
                # Schedule
                'SCHEDULE_TYPE': config.schedule_settings.schedule_type,
                'SCHEDULE_DAY': config.schedule_settings.schedule_day,
                'RUN_TIME': config.schedule_settings.run_time,
                'CRON_SCHEDULE': config.schedule_settings.cron_schedule or "",
                'RUN_ON_STARTUP': config.schedule_settings.run_on_startup,

                # Localization
                'APP_LANGUAGE': config.language.upper(),

                # System
                'HTTP_TIMEOUT': config.http_timeout,
                'LOG_MAX_SIZE_MB': config.logging_settings.max_size_mb,
                'LOG_BACKUP_COUNT': config.logging_settings.backup_count,
            }
            
            # Overlay any explicitly saved preferences
            for k, v in saved_config.items():
                current_config[k] = v
                
            return jsonify(current_config)
        except Exception as e:
            logger.error(f"Error getting config: {e}")
            import traceback
            logger.error(traceback.format_exc())
            return jsonify({'error': str(e)}), 500
            
    elif request.method == 'POST':
        try:
            new_config = request.json
            if not new_config:
                return jsonify({'error': 'No configuration data provided'}), 400
                
            existing_config = load_config_from_file()
            
            for k, v in new_config.items():
                existing_config[k] = v
                
            if save_config_to_file(existing_config):
                # Reload global config
                logger.info("Configuration updated via API, reloading...")
                
                # Re-initialize the app configuration
                from config.settings import load_config_from_env
                config = load_config_from_env()
                
                # Restart the scheduler with the new configuration
                if scheduler:
                    scheduler.shutdown(wait=False)
                
                # Re-init scheduler
                scheduler = init_scheduler()
                
                return jsonify({
                    'success': True,
                    'message': 'Configuration saved successfully'
                })
            else:
                return jsonify({
                    'success': False,
                    'error': 'Failed to save configuration to file'
                }), 500
                
        except Exception as e:
            logger.error(f"Error saving config: {e}")
            import traceback
            logger.error(traceback.format_exc())
            return jsonify({
                'success': False,
                'error': str(e)
            }), 500


@app.route('/api/languages', methods=['GET'])
def get_languages():
    """Get available languages dynamically"""
    from flask import jsonify
    from utils.localization import get_supported_languages
    try:
        return jsonify(get_supported_languages())
    except Exception as e:
        logger.error(f"Error getting languages: {e}")
        return jsonify({'error': str(e)}), 500


@app.route('/api/timezones', methods=['GET'])
def get_timezones():
    """Get IANA timezone identifiers with human-friendly display names"""
    from flask import jsonify
    from utils.timezones import TIMEZONE_NAME_MAP
    try:
        return jsonify(TIMEZONE_NAME_MAP)
    except Exception as e:
        logger.error(f"Error getting timezones: {e}")
        return jsonify({'error': str(e)}), 500


# Initialize/configure the scheduler
def init_scheduler():
    """
    Initialize and configure the scheduler
    
    Args:
        config: Application configuration
        
    Returns:
        Configured scheduler
    """
    global config, scheduler
    scheduler = BackgroundScheduler(job_defaults={
        'misfire_grace_time': 3600
    })
    
    # Add debug ping job if debug is enabled
    if config.logging_settings.debug_mode:
        scheduler.add_job(log_ping, 'interval', minutes=1, id=JOB_ID_DEBUG_PING)
        logger.info("🪲 Debug mode enabled - adding ping job")


    # log cleanup job
    scheduler.add_job(
        lambda: cleanup_log_files(
            config.logging_settings.log_dir,
            config.logging_settings.max_size_mb
        ),
        'interval',
        hours=24,
        id=JOB_ID_LOG_CLEANUP
    )

    
    # Configure main job
    schedule = config.schedule_settings

    if schedule.cron_schedule:
        logger.info(f"Using custom cron schedule: {schedule.cron_schedule}")
        try:
            scheduler.add_job(
                run_main_job,
                CronTrigger.from_crontab(schedule.cron_schedule, timezone=config.timezone),
                id=JOB_ID_MAIN
            )
        except Exception as e:
            logger.error(f"☠️  Invalid cron schedule: {e}")
            # Fall back to default scheduling
            logger.info("⚠️  Falling back to default schedule")
            schedule.cron_schedule = None
            

    # If no custom cron schedule, use schedule type
    if not schedule.cron_schedule:
        if schedule.schedule_type == "DAILY":
            # Daily job at specified time
            logger.info(f"📅  Scheduling DAILY job at {schedule.hour}:{schedule.minute}")
            scheduler.add_job(
                run_main_job, 
                CronTrigger(hour=schedule.hour, minute=schedule.minute, timezone=config.timezone),
                id=JOB_ID_MAIN
            )
        else:
            # Weekly job at specified time and day
            logger.info(f"📅  Scheduling WEEKLY job at {schedule.hour}:{schedule.minute} "
                       f"on day {schedule.schedule_day}")
            scheduler.add_job(
                run_main_job, 
                CronTrigger(
                    day_of_week=schedule.schedule_day, 
                    hour=schedule.hour, 
                    minute=schedule.minute,
                    timezone=config.timezone
                ),
                id=JOB_ID_MAIN
            )
    

    # Run on startup if configured
    if schedule.run_on_startup:
        logger.info("🚀  Running job on startup")
        scheduler.add_job(
            run_main_job,
            'date', 
            run_date=datetime.datetime.now() + datetime.timedelta(seconds=5),
            id='startup_job'
        )
    
    # Start the scheduler
    scheduler.start()
    logger.info("👍  Scheduler started")
    return scheduler

if __name__ == "__main__":
    try:
        # Initialize config and scheduler (which also adds jobs and starts)
        config = init_app()
        scheduler = init_scheduler()

        # Disable all Flask logging
        log = logging.getLogger('werkzeug')
        log.disabled = True

        logger.info("🚀  Starting Flask server...") # Add log message
        app.run(
            host='0.0.0.0',
            port=5000,
            debug=False,
            use_reloader=False,
            threaded=True
        )

    except ValueError as e:
        # Handle configuration validation errors during init_app()
        # Check if logger exists before using it
        if logger:
            logger.error(f"Configuration error: {e}")
        else:
            print(f"Configuration error before logger initialization: {e}")
        sys.exit(1)
    except Exception as e:
        # Catch other potential startup errors
        if logger:
            logger.error(f"An unexpected error occurred during startup: {e}")
            logger.debug(traceback.format_exc())
        else:
            print(f"Unexpected error before logger initialization: {e}")
            print(traceback.format_exc())
        sys.exit(1)