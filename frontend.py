import streamlit as st
import requests

def get_api_response(request_question: str):
    api_url = "http://localhost:8080/send-message"
    headers = {
        "Content-Type": "application/json",
    }
    payload = {
        "question": request_question
    }

    try:
        response = requests.post(api_url, json=payload, headers=headers)
        if response.status_code == 200:
            return response.json()["response"]
        else:
            return {"error": "Failed Response!"}
    except Exception as e:
        return "Error::Something went wrong! " + str(e)


st.title("Assistant for C Level")

question = st.text_input("Ask a question:")

if st.button("Submit"):
    if question:
        with st.spinner("Getting the answer..."):
            answer = get_api_response(question)
            if "Error" in answer:
                st.error(answer)
            else:
                st.write(answer)
    else:
        st.warning("Please enter the question before submitting.")
