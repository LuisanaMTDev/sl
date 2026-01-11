use chrono::Local;
use clap::Parser;
use dotenv::dotenv;
use owo_colors::OwoColorize;
use terminal_size::{Width, terminal_size};

use crate::library::commands::Args;
use crate::library::errors::SLError;
use crate::library::http_client::SLHttpClient;

mod library;

#[tokio::main]
async fn main() -> Result<(), SLError> {
    if let Err(error) = dotenv() {
        return Err(SLError::LoadEnvsFaild(error));
    }

    let sl_http_client = SLHttpClient::new()?;
    let args = Args::parse();

    if args.save && args.topic.is_some() && args.amount_of_cards.is_some() {
        let response = sl_http_client
            .post_lesson(
                args.topic.clone().unwrap(),
                args.amount_of_cards.unwrap(),
                args.calculate_repetitions_dates(),
            )
            .await;

        match response {
            Ok(response_data) => println!("{}", response_data.green().bold()),
            Err(error) => return Err(SLError::RequestToServerFailed(error)),
        };
        if let Some((Width(w), _)) = terminal_size() {
            println!(
                "{:^width$}",
                format!("Today is {}", Local::now().date_naive().bold()),
                width = w as usize
            );
            println!(
                "{:^width$}",
                format!(
                    "Topic \"{}\" has {} amount of cards, and your desired repetitions are {}",
                    args.topic.clone().unwrap(),
                    args.amount_of_cards.unwrap(),
                    args.repetitions_desired
                )
                .bold()
                .yellow(),
                width = w as usize
            );
            args.print_repetitions_dates();
        } else {
            println!(
                "{}",
                format!("Today is {}", Local::now().date_naive().bold()),
            );
            println!(
                "{}",
                format!(
                    "Topic \"{}\" has {} amount of cards, and your desired repetitions are {}",
                    args.topic.clone().unwrap(),
                    args.amount_of_cards.unwrap(),
                    args.repetitions_desired
                )
                .bold()
                .yellow(),
            );
            args.print_repetitions_dates();
        }
    }

    if args.login {
        println!("Start login...");
        let password = match rpassword::prompt_password("Enter password: ") {
            Ok(pw) => pw,
            Err(error) => return Err(SLError::PromptPasswordFaild(error)),
        };
        println!("{}", password);
    }

    Ok(())
}
