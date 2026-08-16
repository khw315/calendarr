import pytest
import os
import sys

# Ensure src module is discoverable
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../src')))

from app import app

@pytest.fixture
def client():
    app.config['TESTING'] = True
    with app.test_client() as client:
        yield client

def test_health_endpoint(client):
    """Test the /health endpoint used by Docker healthcheck."""
    response = client.get('/health')
    assert response.status_code == 200
    data = response.get_json()
    assert 'status' in data
    assert data['status'] == 'healthy'
    assert 'timestamp' in data

def test_timezones_endpoint(client):
    """Test the /api/timezones endpoint."""
    response = client.get('/api/timezones')
    assert response.status_code == 200
    data = response.get_json()
    assert isinstance(data, dict)
    assert len(data) > 0
    # verify some common timezones are there
    assert "UTC" in data or "America/New_York" in data

def test_languages_endpoint(client):
    """Test the /api/languages endpoint returns list of dicts."""
    response = client.get('/api/languages')
    assert response.status_code == 200
    data = response.get_json()
    assert isinstance(data, list)
    assert len(data) > 0
    assert 'code' in data[0]
    assert 'name' in data[0]
