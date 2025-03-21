import streamlit as st
import requests

def get_api_response(request_question: str):
    api_url = "http://localhost:8080/send-message"
    headers = {
        "Content-Type": "application/json",
    }
    payload = {
        "question": request_question + " " + "Jawaban harus dalam bahasa Indonesia, tidak boleh dalam bahasa Inggris. Tidak boleh menjawab diluar topik keuangan"
    }

    try:
        response = requests.post(api_url, json=payload, headers=headers)
        if response.status_code == 200:
            return response.json()["response"]
        else:
            return {"error": "Gagal menanggapi!"}
    except Exception as e:
        return "Error::Terjadi kesalahan! " + str(e)


st.title("Asisten Chief Financial Officer")

question = st.text_input("Pertanyaan:")

if st.button("Submit"):
    if question:
        with st.spinner("Mencari jawaban..."):
            answer = get_api_response(question)
            if "Error" in answer:
                st.error(answer)
            else:
                st.write(answer)
    else:
        st.warning("Harap masukkan pertanyaan sebelum mengirim.")
