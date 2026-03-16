from zoneinfo import available_timezones

# Canonical human-friendly timezone names (DST-safe)
REGION_NAMES = {
    # Global
    "UTC": "Universal Coordinated Time",
    "GMT": "Greenwich Mean Time",

    # Americas
    "ET": "Eastern Time",
    "CT": "Central Time",
    "MT": "Mountain Time",
    "PT": "Pacific Time",
    "AT": "Atlantic Time",
    "AKT": "Alaska Time",
    "HAT": "Hawaii-Aleutian Time",

    # Europe
    "WET": "Western European Time",
    "CET": "Central European Time",
    "EET": "Eastern European Time",
    "MSK": "Moscow Time",

    # Asia
    "IST": "India Standard Time",
    "ICT": "Indochina Time",
    "JST": "Japan Standard Time",
    "KST": "Korea Standard Time",
    "CST_CN": "China Standard Time",
    "GST": "Gulf Time",

    # Australia
    "AET": "Australian Eastern Time",
    "ACT": "Australian Central Time",
    "AWST": "Australian Western Standard Time",

    # Africa
    "WAT": "West Africa Time",
    "CAT": "Central Africa Time",
    "EAT": "East Africa Time",
    "SAST": "South Africa Standard Time",
}

# Prefix / exact mappings (authoritative overrides)
PREFIX_MAP = {
    # UTC / GMT
    "UTC": REGION_NAMES["UTC"],
    "Etc/UTC": REGION_NAMES["UTC"],
    "Etc/GMT": REGION_NAMES["GMT"],
    "Etc/Greenwich": REGION_NAMES["GMT"],

    # ===================== EUROPE =====================
    "Europe/London": REGION_NAMES["GMT"],
    "Europe/Dublin": REGION_NAMES["GMT"],
    "Europe/Lisbon": REGION_NAMES["WET"],
    "Europe/Madeira": REGION_NAMES["WET"],
    "Europe/Canary": REGION_NAMES["WET"],

    "Europe/Paris": REGION_NAMES["CET"],
    "Europe/Berlin": REGION_NAMES["CET"],
    "Europe/Madrid": REGION_NAMES["CET"],
    "Europe/Rome": REGION_NAMES["CET"],
    "Europe/Warsaw": REGION_NAMES["CET"],
    "Europe/Amsterdam": REGION_NAMES["CET"],
    "Europe/Brussels": REGION_NAMES["CET"],
    "Europe/Vienna": REGION_NAMES["CET"],
    "Europe/Zurich": REGION_NAMES["CET"],

    "Europe/Athens": REGION_NAMES["EET"],
    "Europe/Helsinki": REGION_NAMES["EET"],
    "Europe/Bucharest": REGION_NAMES["EET"],
    "Europe/Sofia": REGION_NAMES["EET"],
    "Europe/Tallinn": REGION_NAMES["EET"],

    "Europe/Moscow": REGION_NAMES["MSK"],

    # ===================== AMERICAS =====================
    # Core US/Canada
    "America/New_York": REGION_NAMES["ET"],
    "America/Detroit": REGION_NAMES["ET"],
    "America/Toronto": REGION_NAMES["ET"],
    "America/Chicago": REGION_NAMES["CT"],
    "America/Winnipeg": REGION_NAMES["CT"],
    "America/Denver": REGION_NAMES["MT"],
    "America/Phoenix": REGION_NAMES["MT"],
    "America/Los_Angeles": REGION_NAMES["PT"],
    "America/Vancouver": REGION_NAMES["PT"],

    "America/Anchorage": REGION_NAMES["AKT"],
    "America/Juneau": REGION_NAMES["AKT"],
    "America/Nome": REGION_NAMES["AKT"],
    "America/Honolulu": REGION_NAMES["HAT"],

    # Groups
    "America/Indiana/": REGION_NAMES["ET"],
    "America/North_Dakota/": REGION_NAMES["CT"],
    "America/Argentina/": "Argentina Time",
    "America/Brazil/": "Brazil Time",
    "America/Mexico_": "Mexico Time",
    "America/Caribbean": REGION_NAMES["AT"],

    # South America
    "America/Sao_Paulo": "Brasilia Time",
    "America/Santiago": "Chile Time",
    "America/Bogota": "Colombia Time",
    "America/Lima": "Peru Time",
    "America/Montevideo": "Uruguay Time",

    # ===================== AFRICA =====================
    "Africa/Abidjan": REGION_NAMES["GMT"],
    "Africa/Accra": REGION_NAMES["GMT"],
    "Africa/Lagos": REGION_NAMES["WAT"],
    "Africa/Kinshasa": REGION_NAMES["WAT"],

    "Africa/Harare": REGION_NAMES["CAT"],
    "Africa/Lusaka": REGION_NAMES["CAT"],
    "Africa/Maputo": REGION_NAMES["CAT"],

    "Africa/Nairobi": REGION_NAMES["EAT"],
    "Africa/Kampala": REGION_NAMES["EAT"],
    "Africa/Dar_es_Salaam": REGION_NAMES["EAT"],

    "Africa/Johannesburg": REGION_NAMES["SAST"],
    "Africa/Maseru": REGION_NAMES["SAST"],
    "Africa/Mbabane": REGION_NAMES["SAST"],

    # ===================== ASIA =====================
    "Asia/Kolkata": REGION_NAMES["IST"],
    "Asia/Colombo": REGION_NAMES["IST"],

    "Asia/Bangkok": REGION_NAMES["ICT"],
    "Asia/Ho_Chi_Minh": REGION_NAMES["ICT"],
    "Asia/Vientiane": REGION_NAMES["ICT"],
    "Asia/Phnom_Penh": REGION_NAMES["ICT"],

    "Asia/Tokyo": REGION_NAMES["JST"],
    "Asia/Seoul": REGION_NAMES["KST"],
    "Asia/Shanghai": REGION_NAMES["CST_CN"],

    "Asia/Dubai": REGION_NAMES["GST"],
    "Asia/Qatar": REGION_NAMES["GST"],
    "Asia/Riyadh": REGION_NAMES["GST"],

    "Asia/Singapore": "Singapore Time",
    "Asia/Hong_Kong": "Hong Kong Time",
    "Asia/Taipei": "Taipei Time",

    # Indonesia (official zones)
    "Asia/Jakarta": "Western Indonesia Time",
    "Asia/Makassar": "Central Indonesia Time",
    "Asia/Jayapura": "Eastern Indonesia Time",

    # ===================== AUSTRALIA =====================
    "Australia/Sydney": REGION_NAMES["AET"],
    "Australia/Melbourne": REGION_NAMES["AET"],
    "Australia/Brisbane": REGION_NAMES["AET"],
    "Australia/Adelaide": REGION_NAMES["ACT"],
    "Australia/Darwin": REGION_NAMES["ACT"],
    "Australia/Perth": REGION_NAMES["AWST"],

    # ===================== ATLANTIC / INDIAN =====================
    "Atlantic/Reykjavik": REGION_NAMES["GMT"],
    "Atlantic/Canary": REGION_NAMES["WET"],
    "Atlantic/Azores": "Azores Time",

    "Indian/Maldives": "Maldives Time",
    "Indian/Mauritius": "Mauritius Time",
    "Indian/Reunion": "Reunion Time",

    # ===================== PACIFIC =====================
    "Pacific/Auckland": "New Zealand Time",
    "Pacific/Fiji": "Fiji Time",
    "Pacific/Tahiti": "Tahiti Time",
    "Pacific/Guam": "Chamorro Time",
}

def _fallback_name(tz: str) -> str:
    parts = tz.split("/")
    if len(parts) == 1:
        return tz

    city = parts[-1].replace("_", " ")

    # Avoid ugly names like "Indiana/Knox Time"
    if parts[0] in {"America", "Europe", "Asia"}:
        return city

    return f"{city} Time"

def generate_timezone_map() -> dict[str, str]:
    tz_map: dict[str, str] = {}

    for tz in sorted(available_timezones()):
        matched = False

        for key, value in PREFIX_MAP.items():
            if tz == key:
                tz_map[tz] = value
                matched = True
                break
            if key.endswith("/") and tz.startswith(key):
                tz_map[tz] = value
                matched = True
                break

        if not matched:
            tz_map[tz] = _fallback_name(tz)

    return tz_map

# Export
TIMEZONE_NAME_MAP = generate_timezone_map()
