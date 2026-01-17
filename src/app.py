#!/usr/bin/env python3
# src/app.py

import sys
import io
from contextlib import redirect_stdout, redirect_stderr
import traceback
import datetime
import os
import logging

from flask import Flask
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
        
        # logger.info(f"Job output: {output.getvalue()}")
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
    from flask import send_from_directory
    import os
    public_dir = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'public')
    return send_from_directory(public_dir, 'index.html')

@app.route('/<path:path>')
def serve_static(path):
    """Serve static files"""
    from flask import send_from_directory
    import os
    public_dir = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'public')
    return send_from_directory(public_dir, path)

@app.route('/health')
def health():
    """Health check endpoint for Docker"""
    return {"status": "healthy", "timestamp": datetime.datetime.now().isoformat()}

@app.route('/api/events')
def get_events():
    """Get upcoming events"""
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
            
            # Combine TV and movie events
            for event_item in day.tv_events:
                # Always skip past events in the Web UI
                if event_item.is_past:
                    continue
                    
                day_events.append({
                    'title': event_item.title if hasattr(event_item, 'title') else event_item.summary,
                    'time': event_item.time_str if hasattr(event_item, 'time_str') else None,
                    'type': 'tv'
                })
            
            for event_item in day.movie_events:
                # Always skip past events in the Web UI
                if event_item.is_past:
                    continue
                    
                day_events.append({
                    'title': event_item.title if hasattr(event_item, 'title') else event_item.summary,
                    'time': event_item.time_str if hasattr(event_item, 'time_str') else None,
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



@app.route('/api/past-events')
def get_past_events():
    """Get today's past events"""
    from flask import jsonify
    try:
        # Get only today's events
        from datetime import datetime, timedelta
        import pytz
        
        tz = config.timezone_obj
        now = datetime.now(tz)
        start_date = now.replace(hour=0, minute=0, second=0, microsecond=0)
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
            # Check TV events
            for event_item in day.tv_events:
                if event_item.is_past:
                    past_events.append({
                        'title': event_item.title if hasattr(event_item, 'title') else event_item.summary,
                        'time': event_item.time_str if hasattr(event_item, 'time_str') else None,
                        'type': 'tv'
                    })
            
            # Check movie events
            for event_item in day.movie_events:
                if event_item.is_past:
                    past_events.append({
                        'title': event_item.title if hasattr(event_item, 'title') else event_item.summary,
                        'time': event_item.time_str if hasattr(event_item, 'time_str') else None,
                        'type': 'movie'
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
                CronTrigger.from_crontab(schedule.cron_schedule),
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
                CronTrigger(hour=schedule.hour, minute=schedule.minute),
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
                    minute=schedule.minute
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