import os
import psycopg2
from psycopg2.extras import RealDictCursor
from dotenv import load_dotenv
from pinecone import Pinecone
from openai import OpenAI
import json

# Load environment variables from .env file
load_dotenv()

# OpenAI setup
openai_client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

def get_embedding(text, model="text-embedding-3-small"):
    text = text.replace("\n", " ")
    return openai_client.embeddings.create(input=[text], model=model).data[0].embedding

# Pinecone setup
pinecone_api_key = os.getenv("PINECONE_API_KEY")
pinecone = Pinecone(api_key=pinecone_api_key)

# Connect to your Pinecone index
index_name = 'fusor-threads'
index = pinecone.Index(index_name)

# Connect to PostgreSQL database
conn = psycopg2.connect(
    dbname=os.getenv("DB_NAME"),
    user=os.getenv("DB_USER"),
    password=os.getenv("DB_PASSWORD"),
    host=os.getenv("DB_HOST")
)

def fetch_data_from_db():
    with conn.cursor(cursor_factory=RealDictCursor) as cur:
        cur.execute("SELECT * FROM threads;")
        rows = cur.fetchall()
    return rows

def generate_embeddings(contents):
    return [get_embedding(content) for content in contents]

def insert_into_pinecone(rows):
    embeddings = generate_embeddings([row["content"] for row in rows])
    
    # Prepare the data for Pinecone upsert
    pinecone_data = [
        {
            "id": row["post_id"],
            "values": emb,
            "metadata": { 
                "row_id": row["id"],
                "thread_id": row["thread_id"],
                "content": row["content"],
                "author": row["author"],
                "post_time": row["post_time"].isoformat()  # Convert datetime to string in ISO format
            }
        }
        for row, emb in zip(rows, embeddings)
    ]
    
    # Upsert data into Pinecone
    index.upsert(vectors=pinecone_data)

    print("Data upserted into Pinecone successfully.")

def main():
    rows = fetch_data_from_db()
    rows = [rows[0]]
    if rows:
        insert_into_pinecone(rows)
    else:
        print("No data fetched from the database.")

if __name__ == "__main__":
    main()
