"""Utility functions for localized text and formatting."""

from __future__ import annotations

import random
from typing import Dict, Iterable, Mapping, Optional

SUPPORTED_LANGUAGES = {"EN", "ID", "KO"}
DEFAULT_LANGUAGE = "EN"

# Translation data structured to support simple lookups and count-aware labels.
TRANSLATIONS: Dict[str, Dict[str, object]] = {
    "EN": {
        "movies_heading": "MOVIES",
        "mention_instruction": "If you'd like to be notified when new content is available, join this role!",
        "timezone_message": "All times shown in {timezone}",
        "no_new_releases_messages": [
            "Schedule empty; enjoy some free time.",
            "No new releases today—perfect moment to relax.",
            "Nothing on the calendar right now."
        ],
        "no_day_content_messages": [
            "No releases scheduled for this day.",
            "Nothing new is planned for today.",
            "The schedule is clear for this day."
        ],
        "subheader_templates": {
            "tv": " 📺  {count} {label}",
            "movie": " 🎬  {count} {label}",
            "premiere": " 🎉  {count} {label}"
        },
        "subheader_labels": {
            "tv": {"singular": "all-new episode", "plural": "all-new episodes"},
            "movie": {"singular": "movie release", "plural": "movie releases"},
            "premiere": {"singular": "season premiere", "plural": "season premieres"}
        },
        "list_join": {"two": " and ", "last": ", and "},
        "comma_separator": ", "
    },
    "ID": {
        "movies_heading": "FILM",
        "mention_instruction": "Jika ingin diberi tahu saat ada konten baru, bergabunglah dengan role ini!",
        "timezone_message": "Semua waktu ditampilkan dalam {timezone}",
        "no_new_releases_messages": [
            "Jadwal kosong; nikmati waktu luang.",
            "Tidak ada rilis baru hari ini—saat yang pas untuk bersantai.",
            "Belum ada acara di kalender sekarang."
        ],
        "no_day_content_messages": [
            "Tidak ada jadwal untuk hari ini.",
            "Hari ini belum ada rilis baru.",
            "Kalender hari ini masih kosong."
        ],
        "subheader_templates": {
            "tv": " 📺  {count} {label}",
            "movie": " 🎬  {count} {label}",
            "premiere": " 🎉  {count} {label}"
        },
        "subheader_labels": {
            "tv": "episode terbaru",
            "movie": "rilis film",
            "premiere": "penayangan perdana musim"
        },
        "list_join": {"two": " dan ", "last": ", dan "},
        "comma_separator": ", "
    },
    "KO": {
        "movies_heading": "영화",
        "mention_instruction": "새로운 콘텐츠가 올라오면 알림을 받고 싶다면 이 역할에 참여하세요!",
        "timezone_message": "{timezone} 기준으로 시간이 표시됩니다",
        "no_new_releases_messages": [
            "새로운 일정이 없습니다. 잠시 쉬어가세요.",
            "이번에는 공개 예정이 없어요. 휴식을 즐겨보세요.",
            "캘린더가 비어 있습니다."
        ],
        "no_day_content_messages": [
            "오늘은 예정된 공개가 없습니다.",
            "이 날에는 새로운 일정이 없어요.",
            "등록된 일정이 없습니다."
        ],
        "subheader_templates": {
            "tv": " 📺  {label} {count}편",
            "movie": " 🎬  {label} {count}편",
            "premiere": " 🎉  {label} {count}회"
        },
        "subheader_labels": {
            "tv": "최신 에피소드",
            "movie": "새 영화",
            "premiere": "시즌 첫 방송"
        },
        "list_join": {"two": " 및 ", "last": ", 및 "},
        "comma_separator": ", "
    }
}


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


def get_mention_instruction(language: str) -> str:
    """Return the localized instruction for role mentions."""
    instruction = get_text(language, "mention_instruction")
    if instruction:
        return instruction
    return get_text(DEFAULT_LANGUAGE, "mention_instruction") or ""
