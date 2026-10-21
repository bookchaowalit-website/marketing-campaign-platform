"""verify_token accepts only unexpired HS256 tokens signed with JWT_SECRET."""
import time

import jwt
import pytest

from app.core.auth import verify_token
from app.core.config import settings


def _token(payload, key=None, algorithm="HS256"):
    return jwt.encode(payload, key or settings.JWT_SECRET, algorithm=algorithm)


def test_valid_token_returns_payload():
    payload = verify_token(_token({"user_id": "u1", "exp": int(time.time()) + 60}))
    assert payload["user_id"] == "u1"


def test_wrong_secret_is_rejected():
    assert verify_token(_token({"user_id": "u1"}, key="another-secret-of-32-bytes-length!")) is None


def test_expired_token_is_rejected():
    assert verify_token(_token({"user_id": "u1", "exp": int(time.time()) - 10})) is None


@pytest.mark.parametrize("algorithm", ["HS384", "HS512"])
def test_other_algorithms_are_rejected(algorithm):
    assert verify_token(_token({"user_id": "u1"}, algorithm=algorithm)) is None


def test_unsigned_token_is_rejected():
    assert verify_token(jwt.encode({"user_id": "u1"}, None, algorithm="none")) is None


def test_garbage_is_rejected():
    assert verify_token("not-a-jwt") is None
