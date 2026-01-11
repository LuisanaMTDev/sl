use thiserror::Error;

#[derive(Error, Debug)]
pub enum SLError {
    #[error("ERROR while loading .env file: {0}")]
    LoadEnvsFaild(dotenv::Error),
    #[error("ERROR request failed: {0}")]
    RequestToServerFailed(reqwest::Error),
    #[error("ERROR while readding env var: {0}")]
    ReadEnvVarFailed(std::env::VarError),
    #[error("ERROR while promting password: {0}")]
    PromptPasswordFaild(std::io::Error),
}
