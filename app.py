import streamlit as st
from ai import generate_query_vector, query_pinecone, generate_response_from_gpt4

def chat():
    st.title('Fusor.net Chatbot')
    
    user_question = st.text_input("Ask your question:")
    
    if user_question:
        # Generate a query vector from the user's input
        query_vector = generate_query_vector(user_question)
        
        # Query Pinecone with the embedding
        pinecone_response = query_pinecone(query_vector)

        if pinecone_response:
            # Construct the context from Pinecone's response (retrieved documents)
            context = "\n\n".join([match["content"] for match in pinecone_response])
            context = context.strip()
            
            # Prepare the query with the context for GPT-4-turbo
            query = f"""
            Use the below information to answer the subsequent question. If the answer cannot be found, write "I don't know."
            Information:
            \"\"\"
            {context}
            \"\"\"

            Question: {user_question}
            """
            
            # Prepare messages for GPT-4-turbo
            messages = [
                {'role': 'system', 'content': 'You answer questions based on the provided information.'},
                {'role': 'user', 'content': query},
            ]
            
            # Generate response from GPT-4-turbo
            gpt4_response = generate_response_from_gpt4(messages)
            
            # Display the response
            st.write(gpt4_response.choices[0].message.content)
        else:
            st.write("Sorry, I couldn't find relevant information to answer your question.")

if __name__ == "__main__":
    chat()
