from setuptools import setup, find_packages

setup(
    name="data-ingestion",
    version="0.1.0",
    description="Data ingestion service with database migrations",
    packages=find_packages(where="src"),
    package_dir={"": "src"},
    python_requires=">=3.8",
    install_requires=[
        "alembic>=1.13.0",
        "sqlalchemy>=2.0.0",
        "psycopg2-binary>=2.9.0",
        "python-dotenv>=1.0.0",
    ],
)