import os
from tqdm import tqdm
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
    if text is None or len(text) == 0:
        print("Empty text")
        return None
    try:
        text = text.replace("\n", " ")
        response = openai_client.embeddings.create(input=[text], model=model)
        return response.data[0].embedding
    except Exception as e:
        print(f"An error occurred: {e}")
        print(f"Failed text: {text}")
        return None 


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
        cur.execute("SELECT * FROM threads ORDER BY id;")
        rows = cur.fetchall()
    return rows

def generate_embeddings(contents):
    return [get_embedding(content) for content in contents]

def insert_into_pinecone(rows):
    embeddings = generate_embeddings([row["content"] for row in rows])
    
    # Filter out None values and keep track of the corresponding rows
    valid_embeddings = [(row, emb) for row, emb in zip(rows, embeddings) if emb is not None]
    
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
        for row, emb in valid_embeddings
    ]
    
    if pinecone_data:
        # Upsert data into Pinecone
        index.upsert(vectors=pinecone_data)
        print(f"Batch of {len(pinecone_data)} items upserted into Pinecone successfully.")
    else:
        print("No valid data to upsert into Pinecone.")

def main():
    rows = fetch_data_from_db()
    batch_size = 10

    if rows:
        for i in tqdm(range(0, len(rows), batch_size)):
            batch_rows = rows[i:i+batch_size]
            insert_into_pinecone(batch_rows)
    else:
        print("No data fetched from the database.")

if __name__ == "__main__":
    main()