"""Utility functions for localized text and formatting."""

from __future__ import annotations

import random
from datetime import date, datetime
from typing import Dict, Iterable, Mapping, Optional

SUPPORTED_LANGUAGES = {"EN", "ID", "KO", "JA"}
DEFAULT_LANGUAGE = "EN"

# Translation data structured to support simple lookups and count-aware labels.
import json
from pathlib import Path

# Load translations from JSON file
def _load_translations() -> Dict[str, Dict[str, object]]:
    try:
        data_path = Path(__file__).parent.parent / "data" / "locales.json"
        with open(data_path, "r", encoding="utf-8") as f:
            return json.load(f)
    except Exception as e:
        print(f"Error loading translations: {e}")
        return {}

TRANSLATIONS: Dict[str, Dict[str, object]] = _load_translations()


def normalize_language(language: Optional[str]) -> str:
    """Normalize a language value from environment variables."""
    if not language:
        return DEFAULT_LANGUAGE
    normalized = language.strip().upper()
    return normalized if normalized in SUPPORTED_LANGUAGES else DEFAULT_LANGUAGE


def _get_translation(language: str) -> Dict[str, object]:
    return TRANSLATIONS.get(language, TRANSLATIONS[DEFAULT_LANGUAGE])


def get_text(language: str, key: str) -> Optional[str]:
    """Return a simple string translation for a given key."""
    value = _get_translation(language).get(key)
    return value if isinstance(value, str) else None


def get_random_message(language: str, message_key: str) -> str:
    """Return a random message for a given key, falling back to defaults."""
    translation = _get_translation(language)
    messages = translation.get(f"{message_key}_messages")
    if not isinstance(messages, Iterable):
        messages = _get_translation(DEFAULT_LANGUAGE).get(f"{message_key}_messages", [])
    messages_list = list(messages) if messages else []
    if not messages_list:
        # Provide safe fallbacks when no translations exist
        fallback_map = {
            "no_new_releases": "No new releases to share.",
            "no_day_content": "No releases scheduled for this day."
        }
        return fallback_map.get(message_key, "")
    return random.choice(messages_list)


def _resolve_label(label_config: object, count: int) -> str:
    if isinstance(label_config, Mapping):
        if count == 1 and "singular" in label_config:
            return str(label_config["singular"])
        if count != 1 and "plural" in label_config:
            return str(label_config["plural"])
        if "other" in label_config:
            return str(label_config["other"])
        # Fallback to any available value
        for value in label_config.values():
            return str(value)
        return ""
    return str(label_config) if label_config is not None else ""


def format_subheader_section(language: str, section: str, count: int) -> str:
    """Return a localized subheader fragment for the given section and count."""
    translation = _get_translation(language)
    templates = translation.get("subheader_templates", {})
    labels = translation.get("subheader_labels", {})
    template = templates.get(section)
    label_config = labels.get(section)

    if not isinstance(template, str):
        # Fallback to default language template
        default_translation = _get_translation(DEFAULT_LANGUAGE)
        template = default_translation.get("subheader_templates", {}).get(section, " {count}")
        label_config = default_translation.get("subheader_labels", {}).get(section, "")

    label_text = _resolve_label(label_config, count)
    return template.format(count=count, label=label_text).strip()


def join_with_conjunction(language: str, parts: Iterable[str]) -> str:
    """Join parts using localized conjunction rules."""
    parts_list = [part for part in parts if part]
    if not parts_list:
        return ""
    if len(parts_list) == 1:
        return parts_list[0]

    translation = _get_translation(language)
    joiners = translation.get("list_join", {})
    comma_separator = translation.get("comma_separator", ", ")

    two_join = joiners.get("two", " and ")
    last_join = joiners.get("last", ", and ")

    if len(parts_list) == 2:
        return f"{parts_list[0]}{two_join}{parts_list[1]}"

    leading = comma_separator.join(parts_list[:-1])
    return f"{leading}{last_join}{parts_list[-1]}"


def get_timezone_message(language: str, timezone_name: str) -> str:
    """Return the localized timezone information message."""
    template = get_text(language, "timezone_message")
    if not template:
        template = get_text(DEFAULT_LANGUAGE, "timezone_message") or "All times shown in {timezone}"
    return template.format(timezone=timezone_name)


def get_movies_heading(language: str) -> str:
    """Return the localized heading for movie sections."""
    heading = get_text(language, "movies_heading")
    if heading:
        return heading
    return get_text(DEFAULT_LANGUAGE, "movies_heading") or "MOVIES"


def get_header_text(language: str) -> str:
    """Return the localized header text."""
    header = get_text(language, "header_text")
    if header:
        return header
    return get_text(DEFAULT_LANGUAGE, "header_text") or "New Releases"


def get_mention_instruction(language: str) -> str:
    """Return the localized instruction for role mentions."""
    instruction = get_text(language, "mention_instruction")
    if instruction:
        return instruction
    return get_text(DEFAULT_LANGUAGE, "mention_instruction") or ""


def get_day_name(language: str, weekday: int) -> str:
    """
    Return the localized full day name for a given weekday (0=Monday, 6=Sunday).
    """
    translation = _get_translation(language)
    days = translation.get("days", {})
    
    # Fallback to default language if specific day is missing
    key = str(weekday)
    if key not in days:
        days = _get_translation(DEFAULT_LANGUAGE).get("days", {})
        
    return days.get(key, "")


def get_short_day_name(language: str, weekday: int) -> str:
    """
    Return the localized short day name for a given weekday (0=Monday, 6=Sunday).
    """
    translation = _get_translation(language)
    short_days = translation.get("short_days", {})
    
    # Fallback to default language if specific day is missing
    key = str(weekday)
    if key not in short_days:
        short_days = _get_translation(DEFAULT_LANGUAGE).get("short_days", {})
        
    return short_days.get(key, "")


def get_month_name(language: str, month: int) -> str:
    """
    Return the localized full month name for a given month (1=January, 12=December).
    """
    translation = _get_translation(language)
    months = translation.get("months", {})
    
    # Fallback to default language if specific month is missing
    key = str(month)
    if key not in months:
        months = _get_translation(DEFAULT_LANGUAGE).get("months", {})
        
    return months.get(key, "")


def get_short_month_name(language: str, month: int) -> str:
    """
    Return the localized short month name for a given month (1=January, 12=December).
    """
    translation = _get_translation(language)
    short_months = translation.get("short_months", {})
    
    # Fallback to default language if specific month is missing
    key = str(month)
    if key not in short_months:
        short_months = _get_translation(DEFAULT_LANGUAGE).get("short_months", {})
        
    return short_months.get(key, "")

def format_date(date_obj: date | datetime, language: str) -> str:
    """
    Format a date using the localized 'date_header' template.
    Supported placeholders: {year}, {month}, {day}, {day_name}, {short_day_name}, {month_name}, {short_month_name}
    """
    template = get_text(language, "date_header")
    # Fallback if specific language template is missing
    if not template:
        template = get_text(DEFAULT_LANGUAGE, "date_header") or "{day_name}, {short_month_name} {day}"

    # Extract components
    year = str(date_obj.year)
    month = f"{date_obj.month:02d}"
    day = f"{date_obj.day:02d}"
    
    # Get localized names
    day_name = get_day_name(language, date_obj.weekday())
    short_day_name = get_short_day_name(language, date_obj.weekday())
    month_name = get_month_name(language, date_obj.month)
    short_month_name = get_short_month_name(language, date_obj.month)

    return template.format(
        year=year,
        month=month,
        day=day,
        day_name=day_name,
        short_day_name=short_day_name,
        month_name=month_name,
        short_month_name=short_month_name
    )
