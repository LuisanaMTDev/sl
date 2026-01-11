use chrono::{Local, TimeDelta};
use clap::Parser;
use owo_colors::OwoColorize;
use terminal_size::{Width, terminal_size};

/// CLI to use LuisanaMTDev's spaced learning api from the terminal
#[derive(Debug, Parser)]
#[command(version, about, long_about = None)]
pub struct Args {
    pub topic: Option<String>,
    pub amount_of_cards: Option<i16>,
    #[arg(long = "rd", default_value = "7")]
    pub repetitions_desired: i8,
    #[arg(short = 'a')]
    /// Save topic in database using API.
    pub save: bool,
    #[arg(long = "login")]
    /// Generate API Key
    pub login: bool,
}

impl Args {
    pub fn calculate_repetitions_dates(&self) -> Vec<String> {
        let mut repetitions_dates: Vec<String> = vec![];
        let mut last_repetition_date = Local::now().date_naive();

        for r in 0..=self.repetitions_desired {
            if r == self.repetitions_desired {
                break;
            }

            let new_repetition_date = last_repetition_date + TimeDelta::days(r.into());

            repetitions_dates.push(new_repetition_date.to_string());

            last_repetition_date = new_repetition_date;
        }

        repetitions_dates
    }

    pub fn print_repetitions_dates(&self) {
        let repetitions_dates = self.calculate_repetitions_dates();

        for (r, repetition_date) in repetitions_dates.iter().enumerate() {
            if let Some((Width(w), _)) = terminal_size() {
                println!(
                    "{:^width$}",
                    format!(
                        "⦿ Repetition {} should be on {}",
                        (r + 1).bold(),
                        repetition_date.bold()
                    ),
                    width = w as usize
                );
            } else {
                println!(
                    "{}",
                    format!(
                        "⦿ Repetition {} should be on {}",
                        (r + 1).bold(),
                        repetition_date.bold()
                    ),
                );
            }
        }
    }
}
