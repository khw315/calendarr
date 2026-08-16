import importlib.resources
from zoneinfo import available_timezones

# ---------------------------------------------------------------------------
# Canonical zone set — read directly from zone1970.tab in the tzdata package.
# zone1970.tab lists ONLY canonical (non-link) IANA zones, so this
# automatically excludes every deprecated alias (Asia/Calcutta, Cuba, etc.)
# and updates whenever the tzdata package is updated.
# ---------------------------------------------------------------------------
def _load_canonical_zones() -> frozenset[str]:
    try:
        import tzdata  # noqa: F401
        pkg = importlib.resources.files('tzdata')
        raw = (pkg / 'zoneinfo' / 'zone1970.tab').read_bytes().decode('utf-8')
        zones = {
            line.split('\t')[2]
            for line in raw.splitlines()
            if line and not line.startswith('#')
        }
        return frozenset(zones) | {'UTC', 'GMT'}
    except Exception:
        # Fallback: use all available zones (no filtering)
        return frozenset(available_timezones())


CANONICAL_ZONES = _load_canonical_zones()


# ---------------------------------------------------------------------------
# Human-friendly region name constants (DST-safe)
# ---------------------------------------------------------------------------
REGION_NAMES = {
    # Global
    "UTC": "Universal Coordinated Time",
    "GMT": "Greenwich Mean Time",

    # Americas
    "ET":  "Eastern Time",
    "CT":  "Central Time",
    "MT":  "Mountain Time",
    "PT":  "Pacific Time",
    "AT":  "Atlantic Time",
    "AKT": "Alaska Time",
    "HAT": "Hawaii-Aleutian Time",

    # Europe
    "WET": "Western European Time",
    "CET": "Central European Time",
    "EET": "Eastern European Time",
    "MSK": "Moscow Time",

    # Asia
    "IST":    "India Standard Time",
    "ICT":    "Indochina Time",
    "JST":    "Japan Standard Time",
    "KST":    "Korea Standard Time",
    "CST_CN": "China Standard Time",
    "GST":    "Gulf Standard Time",

    # Australia
    "AET":  "Australian Eastern Time",
    "ACT":  "Australian Central Time",
    "AWST": "Australian Western Standard Time",

    # Africa
    "WAT":  "West Africa Time",
    "CAT":  "Central Africa Time",
    "EAT":  "East Africa Time",
    "SAST": "South Africa Standard Time",
}


# ---------------------------------------------------------------------------
# Authoritative display-name overrides (exact match or trailing-/ prefix)
# ---------------------------------------------------------------------------
PREFIX_MAP: dict[str, str] = {
    # UTC / GMT
    "UTC": REGION_NAMES["UTC"],
    "GMT": REGION_NAMES["GMT"],

    # ===================== EUROPE =====================
    "Europe/London":     REGION_NAMES["GMT"],
    "Europe/Dublin":     REGION_NAMES["GMT"],
    "Europe/Lisbon":     REGION_NAMES["WET"],

    "Europe/Paris":      REGION_NAMES["CET"],
    "Europe/Berlin":     REGION_NAMES["CET"],
    "Europe/Madrid":     REGION_NAMES["CET"],
    "Europe/Rome":       REGION_NAMES["CET"],
    "Europe/Warsaw":     REGION_NAMES["CET"],
    "Europe/Amsterdam":  REGION_NAMES["CET"],
    "Europe/Brussels":   REGION_NAMES["CET"],
    "Europe/Vienna":     REGION_NAMES["CET"],
    "Europe/Zurich":     REGION_NAMES["CET"],

    "Europe/Athens":     REGION_NAMES["EET"],
    "Europe/Helsinki":   REGION_NAMES["EET"],
    "Europe/Bucharest":  REGION_NAMES["EET"],
    "Europe/Sofia":      REGION_NAMES["EET"],
    "Europe/Tallinn":    REGION_NAMES["EET"],

    "Europe/Moscow":     REGION_NAMES["MSK"],

    # ===================== AMERICAS =====================
    "America/New_York":  REGION_NAMES["ET"],
    "America/Detroit":   REGION_NAMES["ET"],
    "America/Toronto":   REGION_NAMES["ET"],

    "America/Chicago":   REGION_NAMES["CT"],
    "America/Winnipeg":  REGION_NAMES["CT"],

    "America/Denver":    REGION_NAMES["MT"],
    "America/Phoenix":   REGION_NAMES["MT"],

    "America/Los_Angeles": REGION_NAMES["PT"],
    "America/Vancouver":   REGION_NAMES["PT"],

    "America/Anchorage": REGION_NAMES["AKT"],
    "America/Juneau":    REGION_NAMES["AKT"],
    "America/Nome":      REGION_NAMES["AKT"],

    "America/Honolulu":  REGION_NAMES["HAT"],

    # Indiana — ET for most, CT for Knox and Tell_City
    "America/Indiana/Indianapolis": REGION_NAMES["ET"],
    "America/Indiana/Marengo":      REGION_NAMES["ET"],
    "America/Indiana/Petersburg":   REGION_NAMES["ET"],
    "America/Indiana/Vevay":        REGION_NAMES["ET"],
    "America/Indiana/Vincennes":    REGION_NAMES["ET"],
    "America/Indiana/Winamac":      REGION_NAMES["ET"],
    "America/Indiana/Knox":         REGION_NAMES["CT"],
    "America/Indiana/Tell_City":    REGION_NAMES["CT"],

    # Kentucky — ET
    "America/Kentucky/Louisville":  REGION_NAMES["ET"],
    "America/Kentucky/Monticello":  REGION_NAMES["ET"],

    # Jamaica observes EST year-round (no DST)
    "America/Jamaica":              REGION_NAMES["ET"],

    # North Dakota — all CT
    "America/North_Dakota/": REGION_NAMES["CT"],

    # South America
    "America/Argentina/": "Argentina Time",
    "America/Sao_Paulo":  "Brasilia Time",
    "America/Santiago":   "Chile Time",
    "America/Bogota":     "Colombia Time",
    "America/Lima":       "Peru Time",
    "America/Montevideo": "Uruguay Time",

    # ===================== AFRICA =====================
    "Africa/Abidjan": REGION_NAMES["GMT"],
    "Africa/Accra":   REGION_NAMES["GMT"],

    "Africa/Lagos":    REGION_NAMES["WAT"],
    "Africa/Kinshasa": REGION_NAMES["WAT"],

    "Africa/Harare":  REGION_NAMES["CAT"],
    "Africa/Lusaka":  REGION_NAMES["CAT"],
    "Africa/Maputo":  REGION_NAMES["CAT"],

    "Africa/Nairobi":       REGION_NAMES["EAT"],
    "Africa/Kampala":       REGION_NAMES["EAT"],
    "Africa/Dar_es_Salaam": REGION_NAMES["EAT"],

    "Africa/Johannesburg": REGION_NAMES["SAST"],
    "Africa/Maseru":       REGION_NAMES["SAST"],
    "Africa/Mbabane":      REGION_NAMES["SAST"],

    # ===================== ASIA =====================
    "Asia/Kolkata": REGION_NAMES["IST"],
    "Asia/Colombo": REGION_NAMES["IST"],

    "Asia/Bangkok":    REGION_NAMES["ICT"],
    "Asia/Ho_Chi_Minh": REGION_NAMES["ICT"],
    "Asia/Vientiane":  REGION_NAMES["ICT"],
    "Asia/Phnom_Penh": REGION_NAMES["ICT"],

    "Asia/Tokyo":    REGION_NAMES["JST"],
    "Asia/Seoul":    REGION_NAMES["KST"],
    "Asia/Shanghai": REGION_NAMES["CST_CN"],

    "Asia/Dubai":  REGION_NAMES["GST"],
    "Asia/Qatar":  REGION_NAMES["GST"],
    "Asia/Riyadh": REGION_NAMES["GST"],

    "Asia/Singapore": "Singapore Time",
    "Asia/Hong_Kong": "Hong Kong Time",
    "Asia/Taipei":    "Taipei Time",

    # Indonesia
    "Asia/Jakarta":  "Western Indonesia Time",
    "Asia/Makassar": "Central Indonesia Time",
    "Asia/Jayapura": "Eastern Indonesia Time",

    # ===================== AUSTRALIA =====================
    "Australia/Sydney":    REGION_NAMES["AET"],
    "Australia/Melbourne": REGION_NAMES["AET"],
    "Australia/Brisbane":  REGION_NAMES["AET"],
    "Australia/Adelaide":  REGION_NAMES["ACT"],
    "Australia/Darwin":    REGION_NAMES["ACT"],
    "Australia/Perth":     REGION_NAMES["AWST"],

    # ===================== ATLANTIC / INDIAN =====================
    "Atlantic/Reykjavik": REGION_NAMES["GMT"],
    "Atlantic/Canary":    REGION_NAMES["WET"],
    "Atlantic/Azores":    "Azores Time",

    "Indian/Maldives":  "Maldives Time",
    "Indian/Mauritius": "Mauritius Time",
    "Indian/Reunion":   "Reunion Time",

    # ===================== PACIFIC =====================
    "Pacific/Auckland": "New Zealand Time",
    "Pacific/Fiji":     "Fiji Time",
    "Pacific/Tahiti":   "Tahiti Time",
    "Pacific/Guam":     "Chamorro Time",
}


def _fallback_name(tz: str) -> str:
    """Generate a consistent 'City Time' label for unmapped canonical zones."""
    parts = tz.split("/")
    city = parts[-1].replace("_", " ")
    return f"{city} Time"


def generate_timezone_map() -> dict[str, str]:
    tz_map: dict[str, str] = {}

    for tz in sorted(CANONICAL_ZONES):
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
