use clap::Parser;

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
/// Command line arguments for the application.
pub struct Args {
    /// The task to perform (e.g., "create", "sync", "inspect", "develop").
    #[arg(short, long)]
    pub task: String,

    /// The name of the module to operate on. Optional.
    #[arg(short, long)]
    pub module: Option<String>,

    /// The type of the module (e.g., "full", "light"). Optional.
    #[arg(short, long)]
    pub module_type: Option<String>,

    /// Flag to indicate a dry run. Optional.
    #[arg(short, long)]
    pub dry_run: Option<bool>,

    /// The type of inspection to perform (default is "all").
    #[arg(long, default_value = "all")]
    pub inspect_type: String,
}
