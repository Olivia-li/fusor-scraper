# your_script.py
import os
from dotenv import load_dotenv
from pinecone import Pinecone
from openai import OpenAI

# Load environment variables from .env file
load_dotenv()

# Pinecone setup
pinecone_api_key = os.getenv("PINECONE_API_KEY")
pinecone = Pinecone(api_key=pinecone_api_key)

# Connect to your Pinecone index
index_name = 'fusor-threads'
index = pinecone.Index(index_name)

# OpenAI setup
openai_client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

def generate_query_vector(text, model="text-embedding-3-small"):
    # Generate an embedding for the provided text
    response = openai_client.embeddings.create(input=text, model=model)
    return response.data[0].embedding

def query_pinecone(query_vector, top_k=10):
    # Query Pinecone with the provided vector
    query_response = index.query(vector=query_vector, top_k=top_k, include_metadata=True)
    # Extract the vector IDs from the query response only if the score is above 0.95
    print(query_response["matches"][0]["score"])
    response = [match["metadata"] for match in query_response["matches"] if match["score"] > 0.45]
    # Return the fetched data
    return response

def generate_response_from_gpt4(messages):
    # Generate a response using OpenAI's GPT-4-turbo based on the previous messages
    response = openai_client.chat.completions.create(
        # model="gpt-4-0125-preview",
        model="gpt-3.5-turbo-1106",
        messages=messages
    )
    return response
