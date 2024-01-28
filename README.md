# Fusor Chatbot

## Overview
Fusor Chatbot is an intelligent chatbot application designed to provide accurate and context-aware answers. The application leverages the power of Pinecone for efficient information retrieval and OpenAI's GPT-4-turbo for generating precise and human-like responses. It's a perfect tool for querying complex datasets and providing information in a conversational manner.

## Features
- **Contextual Information Retrieval**: Utilizes Pinecone's vector database to retrieve relevant information based on user queries.
- **Advanced Natural Language Processing**: Leverages OpenAI's GPT-4-turbo to understand and respond to user queries in a natural and intuitive way.
- **Streamlit Interface**: Provides an easy-to-use web interface built with Streamlit, allowing users to interact with the chatbot seamlessly.

## Installation & Setup

### Prerequisites
- Python 3.x
- Poetry for dependency management

### Installation Steps
1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/fusor-chatbot.git
   cd fusor-chatbot

2. **Install dependencies:**
   ```bash
   poetry install

3. **Set up environment variables:**
  Create a .env file in the root directory and add the following variables:
   ```.env
     OPENAI_API_KEY=your_openai_api_key
     PINECONE_API_KEY=your_pinecone_api_key
     DB_NAME=your_database_name
     DB_USER=your_database_user
     DB_PASSWORD=your_database_password
     DB_HOST=your_database_host

4. **Run the application:**
   ```bash
     poetry run streamlit run app.py

## Usage
After starting the application, navigate to the local URL provided by Streamlit (`http://localhost:8501`). Enter your question in the input field and receive a response generated by the chatbot.

## Contributing 
Contributions are welcome! If you have suggestions or improvements, feel free to fork the repository and submit a pull request.

## License
MIT License
