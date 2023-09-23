/*
Copyright (c) 2023 Sienna Lloyd

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

use clap::Parser;
use mimalloc::MiMalloc;
use tracing::{event, Level};
use tracing::instrument::WithSubscriber;
use tracing_subscriber::{EnvFilter, prelude::*};

use pleiades::aspects::config::Root;

#[global_allocator]
static GLOBAL: MiMalloc = MiMalloc;

#[tracing::instrument]
pub fn main() {
    let args = Root::parse();

    let mut def_log_level = Level::INFO;
    if args.debug || args.trace {
        if args.debug {
            def_log_level = Level::DEBUG
        } else {
            def_log_level = Level::TRACE
        };
    }

    // configure tracing & logging
    let reg = tracing_subscriber::registry()
        .with(EnvFilter::from_env("PLEIADES_LOG"))
        .with_subscriber(
            tracing_subscriber::fmt()
                .json()
                .with_max_level(def_log_level)
                .finish()
        );
    tracing::dispatcher::set_global_default(reg.dispatcher().clone());

    event!(Level::INFO, ?args);
    event!(Level::DEBUG, ?args);
}
