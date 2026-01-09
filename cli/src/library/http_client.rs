use std::env;

use reqwest;

use crate::library::errors::SLError;

#[derive(Debug)]
pub struct SLHttpClient {
    client: reqwest::Client,
    api_url: String,
}

impl SLHttpClient {
    pub fn new() -> Result<Self, SLError> {
        let api_url_result = env::var("SL_API_URL");
        match api_url_result {
            Ok(api_url) => Ok(Self {
                client: reqwest::Client::new(),
                api_url,
            }),
            Err(error) => Err(SLError::ReadEnvVarFailed(error)),
        }
    }

    pub async fn post_lesson(
        &self,
        topic: String,
        amount_of_cards: i16,
        repetitions_dates: Vec<String>,
    ) -> Result<String, reqwest::Error> {
        // TODO: Add appropriated error handling.
        let response = self
            .client
            .post(format!("{}/lesson/add", self.api_url))
            .header("SL-Client-Type", "SL-CLI")
            .json(&serde_json::json!({"topic": topic,"amount_of_cards": amount_of_cards, "repetitions_dates": repetitions_dates }))
            .send()
            .await;

        let response_body_result = match response {
            Err(error) => return Err(error),
            Ok(resolved_response) => resolved_response.text().await,
        };

        match response_body_result {
            Err(error) => Err(error),
            Ok(response_body) => Ok(response_body),
        }
    }
}
