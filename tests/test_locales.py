import json
from pathlib import Path

def get_locales_dir():
    # calendarr/tests/test_locales.py -> calendarr/src/data/locales
    return Path(__file__).parent.parent / "src" / "data" / "locales"

def test_locales_match_template():
    locales_dir = get_locales_dir()
    template_path = locales_dir / "template.json"
    
    assert template_path.exists(), "template.json must exist"
    
    with open(template_path, "r", encoding="utf-8") as f:
        template_data = json.load(f)
        
    template_keys = set(template_data.keys())
    template_keys.discard("_comment")
    
    # Check all other json files
    for file_path in locales_dir.glob("*.json"):
        if file_path.name == "template.json":
            continue
            
        with open(file_path, "r", encoding="utf-8") as f:
            lang_data = json.load(f)
            
        lang_keys = set(lang_data.keys())
        lang_keys.discard("_comment")
        
        # Check for missing keys
        missing_keys = template_keys - lang_keys
        assert not missing_keys, f"File {file_path.name} is missing keys: {missing_keys}"
        
        # Types check
        for key in template_keys:
            if key in lang_data:
                assert isinstance(lang_data[key], type(template_data[key])), \
                    f"Type mismatch for key '{key}' in {file_path.name}: expected {type(template_data[key]).__name__}, got {type(lang_data[key]).__name__}"
