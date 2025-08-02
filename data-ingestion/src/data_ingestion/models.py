from sqlalchemy import Column, Integer, String, Float, JSON
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()


class Quote(Base):
    __tablename__ = 'quotes'
    
    id = Column(Integer, primary_key=True, autoincrement=True)
    quote = Column(String, nullable=False)
    author = Column(String, nullable=False)
    tags = Column(JSON, nullable=False)
    popularity = Column(Float, nullable=False)
    category = Column(String, nullable=False)