"""Utility functions for localized text and formatting."""

from __future__ import annotations

import random
from typing import Dict, Iterable, Mapping, Optional

SUPPORTED_LANGUAGES = {"EN", "ID", "KO", "JA"}
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
            "Nothing on the calendar right now.",
            "No upcoming releases at the moment.",
            "The schedule is clear for now; take a break.",
            "No new content to share at this time.",
            "Enjoy your day—no releases scheduled.",
            "The calendar is empty; have a restful time.",
            "No new shows or movies today; unwind and enjoy.",
            "Take it easy—no releases planned for now."
        ],
        "no_day_content_messages": [
            "No releases scheduled for this day.",
            "Nothing new is planned for today.",
            "The schedule is clear for this day.",
            "No new content available for today.",
            "Today has no scheduled releases.",
            "The calendar shows no new releases for today.",
            "Enjoy your day—no new content scheduled.",
            "No new shows or movies planned for today.",
            "Take a break—nothing new is scheduled for today.",
            "The schedule is empty for today; relax and enjoy."
        ],
        "subheader_templates": {
            "tv": " {count} {label}",
            "movie": "{count} {label}",
            "premiere": "{count} {label}"
        },
        "subheader_labels": {
            "tv": {"singular": "all-new episode", "plural": "all-new episodes"},
            "movie": {"singular": "movie release", "plural": "movie releases"},
            "premiere": {"singular": "season premiere", "plural": "season premieres"}
        },
        "list_join": {"two": " and ", "last": ", and "},
        "comma_separator": ", ",
        "days": {
            0: "Monday", 1: "Tuesday", 2: "Wednesday", 3: "Thursday", 4: "Friday", 5: "Saturday", 6: "Sunday"
        },
        "short_days": {
             0: "Mon", 1: "Tue", 2: "Wed", 3: "Thu", 4: "Fri", 5: "Sat", 6: "Sun"
        }
    },
    "ID": {
        "movies_heading": "FILM",
        "mention_instruction": "Jika ingin diberi tahu saat ada konten baru, bergabunglah dengan role ini!",
        "timezone_message": "Semua waktu ditampilkan dalam {timezone}",
        "no_new_releases_messages": [
            "Jadwal kosong; nikmati waktu luang.",
            "Tidak ada rilis baru hari ini—saat yang pas untuk bersantai.",
            "Belum ada acara di kalender sekarang.",
            "Tidak ada rilis mendatang saat ini.",
            "Jadwal kosong untuk saat ini; istirahatlah sejenak.",
            "Tidak ada konten baru untuk dibagikan saat ini.",
            "Nikmati harimu—tidak ada rilis yang dijadwalkan.",
            "Kalender kosong; semoga waktu istirahatmu menyenangkan.",
            "Tidak ada acara atau film baru hari ini; bersantailah dan nikmati.",
            "Santai saja—tidak ada rilis yang direncanakan untuk saat ini."
        ],
        "no_day_content_messages": [
            "Tidak ada jadwal untuk hari ini.",
            "Hari ini belum ada rilis baru.",
            "Kalender hari ini masih kosong.",
            "Tidak ada konten baru untuk hari ini.",
            "Hari ini tidak ada rilis yang dijadwalkan.",
            "Kalender menunjukkan tidak ada rilis baru untuk hari ini.",
            "Nikmati harimu—tidak ada konten baru yang dijadwalkan.",
            "Tidak ada acara atau film baru yang direncanakan untuk hari ini.",
            "Istirahatlah—tidak ada yang dijadwalkan untuk hari ini.",
            "Jadwal hari ini kosong; bersantailah dan nikmati."
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
        "comma_separator": ", ",
        "days": {
            0: "Senin", 1: "Selasa", 2: "Rabu", 3: "Kamis", 4: "Jumat", 5: "Sabtu", 6: "Minggu"
        },
        "short_days": {
             0: "Sen", 1: "Sel", 2: "Rab", 3: "Kam", 4: "Jum", 5: "Sab", 6: "Min"
        }
    },
    "KO": {
        "movies_heading": "영화",
        "mention_instruction": "새로운 콘텐츠가 올라오면 알림을 받고 싶다면 이 역할에 참여하세요!",
        "timezone_message": "{timezone} 기준으로 시간이 표시됩니다",
        "no_new_releases_messages": [
            "새로운 일정이 없습니다. 잠시 쉬어가세요.",
            "이번에는 공개 예정이 없어요. 휴식을 즐겨보세요.",
            "캘린더가 비어 있습니다.",
            "현재 예정된 공개가 없습니다.",
            "지금은 일정이 비어 있으니 잠시 쉬어가세요.",
            "공유할 새로운 콘텐츠가 없습니다.",
            "오늘은 예정된 공개가 없습니다. 좋은 하루 보내세요.",
            "캘린더가 비어 있습니다. 편안한 시간 보내세요.",
            "오늘은 새로운 프로그램이나 영화가 없습니다. 편안히 즐기세요.",
            "지금은 예정된 공개가 없으니 여유롭게 쉬어가세요."
        ],
        "no_day_content_messages": [
            "오늘은 예정된 공개가 없습니다.",
            "이 날에는 새로운 일정이 없어요.",
            "등록된 일정이 없습니다.",
            "오늘은 새로운 콘텐츠가 없습니다.",
            "오늘은 예정된 공개가 없습니다.",
            "오늘은 새로운 공개가 없다고 캘린더에 표시되어 있습니다.",
            "오늘은 예정된 새로운 콘텐츠가 없습니다. 좋은 하루 보내세요.",
            "오늘은 새로운 프로그램이나 영화가 없습니다.",
            "오늘은 새로운 일정이 없으니 잠시 쉬어가세요.",
            "오늘은 일정이 비어 있으니 편안히 즐기세요."
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
        "comma_separator": ", ",
        "days": {
            0: "월요일", 1: "화요일", 2: "수요일", 3: "목요일", 4: "금요일", 5: "토요일", 6: "일요일"
        },
        "short_days": {
             0: "월", 1: "화", 2: "수", 3: "목", 4: "금", 5: "토", 6: "일"
        }
    },
    "JA": {
        "movies_heading": "映画",
        "mention_instruction": "新しいコンテンツの通知を受け取りたい場合は、このロールに参加してください！",
        "timezone_message": "時間は{timezone}基準で表示されます",
        "no_new_releases_messages": [
            "スケジュールは空いています。自由な時間をお楽しみください。",
            "今日の新しいリリースはありません。リラックスするのに最適な時間です。",
            "現在カレンダーに予定はありません。",
            "現在、予定されているリリースはありません。",
            "スケジュールは今のところ空白です。休憩しましょう。",
            "現在共有する新しいコンテンツはありません。",
            "良い一日を — 予定されているリリースはありません。",
            "カレンダーは空です。安らかな時間をお過ごしください。",
            "今日は新しい番組や映画はありません。ゆっくりお楽しみください。",
            "気楽に行きましょう — 今のところリリースの予定はありません。"
        ],
        "no_day_content_messages": [
            "この日のリリースの予定はありません。",
            "今日は新しい予定はありません。",
            "この日のスケジュールは空いています。",
            "今日利用可能な新しいコンテンツはありません。",
            "今日の予定されているリリースはありません。",
            "カレンダーには今日の新しいリリースが表示されていません。",
            "良い一日を — 新しいコンテンツの予定はありません。",
            "今日予定されている新しい番組や映画はありません。",
            "休憩しましょう — 今日は何も予定されていません。",
            "今日のスケジュールは空です。リラックスしてお楽しみください。"
        ],
        "subheader_templates": {
            "tv": " 📺  {label} {count}話",
            "movie": " 🎬  {label} {count}本",
            "premiere": " 🎉  {label} {count}回"
        },
        "subheader_labels": {
            "tv": "新エピソード",
            "movie": "映画リリース",
            "premiere": "シーズン初回"
        },
        "list_join": {"two": " と ", "last": "、そして "},
        "comma_separator": "、",
        "days": {
            0: "月曜日", 1: "火曜日", 2: "水曜日", 3: "木曜日", 4: "金曜日", 5: "土曜日", 6: "日曜日"
        },
        "short_days": {
             0: "月", 1: "火", 2: "水", 3: "木", 4: "金", 5: "土", 6: "日"
        }
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


def get_day_name(language: str, weekday: int) -> str:
    """
    Return the localized full day name for a given weekday (0=Monday, 6=Sunday).
    """
    translation = _get_translation(language)
    days = translation.get("days", {})
    
    # Fallback to default language if specific day is missing
    if weekday not in days:
        days = _get_translation(DEFAULT_LANGUAGE).get("days", {})
        
    return days.get(weekday, "")


def get_short_day_name(language: str, weekday: int) -> str:
    """
    Return the localized short day name for a given weekday (0=Monday, 6=Sunday).
    """
    translation = _get_translation(language)
    short_days = translation.get("short_days", {})
    
    # Fallback to default language if specific day is missing
    if weekday not in short_days:
        short_days = _get_translation(DEFAULT_LANGUAGE).get("short_days", {})
        
    return short_days.get(weekday, "")
