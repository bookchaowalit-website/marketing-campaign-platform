.PHONY: run dev test clean migrate generate-proto

PYTHON=python3
PIP=pip3

run:
	$(PYTHON) -m uvicorn app.main:app --host 0.0.0.0 --port 8081

dev:
	$(PYTHON) -m uvicorn app.main:app --reload --host 0.0.0.0 --port 8081

test:
	$(PYTHON) -m pytest tests/ -v --cov=app --cov-report=html

clean:
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type d -name "*.pyc" -exec rm -rf {} +
	find . -type d -name ".pytest_cache" -exec rm -rf {} +
	find . -type d -name ".mypy_cache" -exec rm -rf {} +

install:
	$(PIP) install -r requirements.txt

install-dev:
	$(PIP) install -r requirements-dev.txt

migrate:
	$(PYTHON) -m alembic upgrade head

migrate-down:
	$(PYTHON) -m alembic downgrade -1

migrate-generate:
	$(PYTHON) -m alembic revision --autogenerate -m "$(MESSAGE)"

generate-proto:
	cd ../../packages/grpc-protos && make generate

db-init:
	$(PYTHON) -m alembic upgrade head

seed:
	$(PYTHON) -m app.utils.seeder
