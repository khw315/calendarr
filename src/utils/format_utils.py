#!/usr/bin/env python3
# src/utils/format_utils.py

# TODO: This file needs more debug logging eventually
import traceback
from typing import Dict, Iterable, Optional

from constants import (
    COLOR_PALETTE, PLATFORM_SLACK, TIMEZONE_NAME_MAP,
    PLATFORM_DISCORD,
    DISCORD_BOLD_START, DISCORD_BOLD_END, SLACK_BOLD_START, SLACK_BOLD_END,
    ITALIC_START, ITALIC_END,
    DEFAULT_HEADER
)
from datetime import datetime
import pytz
import logging

from utils.date_utils import get_days_order
from utils.localization import (
    format_subheader_section,
    get_random_message,
    get_timezone_message,
    join_with_conjunction,
    get_day_name,
    get_short_month_name,
    get_header_text
)

logger = logging.getLogger("format_utils")





# def build_content_summary_parts(tv_count: int, movie_count: int, 
#                                premiere_count: int, platform: str = "discord") -> List[str]:
#     """
#     Build the content summary parts with counts and emojis
    
#     Args:
#         tv_count: Number of TV episodes
#         movie_count: Number of movie releases
#         premiere_count: Number of premieres
#         platform: Platform name for formatting
        
#     Returns:
#         List of strings with formatted content parts
#     """
#     # Different platforms need different spacing after emojis
#     emoji_spacing = "  " if platform == "slack" else " "
    
#     parts = []
#     # Add TV shows count if any
#     if tv_count > 0:
#         shows_text = pluralize("episode", tv_count)
#         parts.append(f"📺{emoji_spacing}{tv_count} all-new {shows_text}")
    
#     # Add movies if any
#     if movie_count > 0:
#         movies_text = pluralize("movie release", movie_count)
#         parts.append(f"🎬{emoji_spacing}{tv_count} {movies_text}")
    
#     # Add premieres if any
#     if premiere_count > 0:
#         premiere_text = pluralize("premiere", premiere_count)
#         parts.append(f"🎉{emoji_spacing}{tv_count} season {premiere_text}")
    
#     return parts


# def join_content_parts(parts: List[str], platform: str = "discord") -> str:
#     """
#     Join content parts with appropriate separators and formatting
    
#     Args:
#         parts: List of content summary parts
#         platform: Platform name for formatting
        
#     Returns:
#         Formatted string with all parts joined
#     """
#     if not parts:
#         return NO_NEW_RELEASES_MSG
    
#     # Different bold syntax for different platforms
#     bold_start = "*" if platform == "slack" else "**"
#     bold_end = "*" if platform == "slack" else "**"
    
#     if len(parts) == 1:
#         return f"{bold_start}{parts[0]}{bold_end}"
#     elif len(parts) == 2:
#         return f"{bold_start}{parts[0]} and {parts[1]}{bold_end}"
#     else:
#         # Join all but last with commas, then add the last with "and"
#         return f"{bold_start}{', '.join(parts[:-1])}, and {parts[-1]}{bold_end}"


def format_header_text(custom_header: str, start_date, end_date, 
                      show_date_range: bool, language: str = "EN") -> str:
    """
    Create a formatted header text with optional date range
    
    Args:
        custom_header: Header text
        start_date: Start date
        end_date: End date
        show_date_range: Whether to show date range
        language: Language code for localization
        
    Returns:
        Formatted header text
    """
    # Use localized header if the provided header matches the default
    if custom_header == DEFAULT_HEADER:
        header_text = get_header_text(language)
    else:
        header_text = f"{custom_header}"
    
    if show_date_range:
        # Check if we're in daily mode (start and end date are the same day)
        if start_date.date() == end_date.date():
            # For daily mode, show just the day name and date
            # day_name = start_date.strftime('%A')
            day_name = get_day_name(language, start_date.weekday())
            
            # month_name = start_date.strftime('%b')
            month_name = get_short_month_name(language, start_date.month)
            day_num = start_date.day
            
            # Format: DayName, Mon DD
            # Localization nuance: Asian languages often put Month before Day (MM-DD), Western (Mon DD or DD Mon).
            # For now, let's stick to the existing structure but with localized strings: "Monday, Jan 01"
            # Or if Asian: "Month Day (DayName)"?
            
            if language == "ID":
                # 1 Jan (Senin)
                header_text += f" ({day_num} {month_name}, {day_name})"
            elif language == "KO":
                # 1월 1일 (월요일)
                # Correct suffix for Korean day is '일'
                header_text += f" ({month_name}월 {day_num}일, {day_name})"
            elif language == "JA":
                # 1月 1日 (月曜日)
                header_text += f" ({month_name}月 {day_num}日, {day_name})"
            else:
                 header_text += f" ({day_name}, {month_name} {day_num:02d})"
                 
        else:
            # For weekly mode, show the range as before
            # Apr 01 - Apr 07
            start_month = get_short_month_name(language, start_date.month)
            start_day = start_date.day
            end_month = get_short_month_name(language, end_date.month)
            end_day = end_date.day
            
            if language == "JA":
                 # 4月 1日 - 4月 7日
                 header_text += f" ({start_month}月 {start_day}日 - {end_month}月 {end_day}日)"
            elif language == "KO":
                 # 4월 1일 - 4월 7일
                 # Correct suffix for Korean day is '일'
                 header_text += f" ({start_month}월 {start_day}일 - {end_month}월 {end_day}일)"
            else:
                 # Apr 01 - Apr 07
                 header_text += f" ({start_month} {start_day:02d} - {end_month} {end_day:02d})"
    
    return header_text


def format_subheader_text(tv_count: int, movie_count: int, premiere_count: int,
                          platform: str, language: str) -> str:
    """
    Format the subheader text showing counts of content, applying platform-specific bolding.

    Args:
        tv_count: Number of TV episodes
        movie_count: Number of movie releases
        premiere_count: Number of premieres
        platform: The target platform ('discord' or 'slack')
        language: ISO language code used for localization

    Returns:
        Formatted subheader text with platform-specific bolding (includes trailing newlines)
    """
    bold_start = SLACK_BOLD_START if platform == PLATFORM_SLACK else DISCORD_BOLD_START
    bold_end = SLACK_BOLD_END if platform == PLATFORM_SLACK else DISCORD_BOLD_END

    # Determine if there are any events at all
    if tv_count == 0 and movie_count == 0:
        nothing_new_message = get_random_message(language, "no_new_releases")
        return f"{bold_start}{nothing_new_message}{bold_end}\n\n"

    subheader_sections = []

    if tv_count > 0:
        section_text = format_subheader_section(language, "tv", tv_count)
        subheader_sections.append(f"{bold_start}{section_text}{bold_end}")

    if movie_count > 0:
        section_text = format_subheader_section(language, "movie", movie_count)
        subheader_sections.append(f"{bold_start}{section_text}{bold_end}")

    if premiere_count > 0:
        section_text = format_subheader_section(language, "premiere", premiere_count)
        subheader_sections.append(f"{bold_start}{section_text}{bold_end}")

    subheader = join_with_conjunction(language, subheader_sections)
    return subheader + "\n\n" if subheader else ""


def get_day_colors(platform: str, start_week_on_monday: bool = True) -> Dict:
    """
    Get ROYGBIV color mapping for days of the week
    
    Args:
        platform: Platform name
        start_week_on_monday: Whether week starts on Monday
        
    Returns:
        Dictionary mapping day names to color codes
    """    

    days_order = get_days_order(start_week_on_monday)

    color_order = ["red", "orange", "yellow", "green", "blue", "indigo", "violet"]

    day_colors = {}
    for i, day in enumerate(days_order):
        color_name = color_order[i % len(color_order)]
        day_colors[day] = COLOR_PALETTE[platform.lower()][color_name]
    
    return day_colors


def format_timezone_line(timezone_obj: Optional[pytz.BaseTzInfo], platform: str,
                         language: str) -> str:
    """Format a timezone information line using localized text."""
    if not timezone_obj:
        logger.warning("‼️  No timezone object provided to format_timezone_line.")
        return ""

    tz_display_name = None
    try:
        tz_identifier = timezone_obj.zone  # e.g., "America/Chicago"
        if tz_identifier in TIMEZONE_NAME_MAP:
            tz_display_name = TIMEZONE_NAME_MAP[tz_identifier]
            logger.debug(
                "🍭  Using custom timezone name '%s' for identifier '%s'.",
                tz_display_name,
                tz_identifier,
            )
        else:
            standard_time_sample = datetime(datetime.now().year, 1, 1)
            localized_sample = timezone_obj.localize(standard_time_sample)
            tz_abbr = localized_sample.tzname()
            if tz_abbr:
                tz_display_name = tz_abbr
                logger.debug(
                    "Using standard time abbreviation '%s' for identifier '%s'.",
                    tz_display_name,
                    tz_identifier,
                )
                if "+" in tz_abbr or "-" in tz_abbr or len(tz_abbr) > 5:
                    logger.warning(
                        "‼️  Timezone abbreviation '%s' might be an offset or non-standard. Using it anyway.",
                        tz_abbr,
                    )
            else:
                logger.warning(
                    "⚠️  Could not determine timezone abbreviation for identifier '%s'.",
                    tz_identifier,
                )

    except Exception as e:
        logger.error(f"☠️  Error determining timezone display name: {e}")
        logger.debug(traceback.format_exc())

    if tz_display_name:
        message = get_timezone_message(language, tz_display_name)
        return f"{ITALIC_START}{message}{ITALIC_END}"

    return ""
