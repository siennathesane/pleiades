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

use std::sync::atomic::Ordering::{Relaxed, Release};
use std::sync::atomic::{AtomicBool, Ordering};
use std::time;

/// The core Hybrid Logical Clock struct.
#[derive(Debug)]
pub struct HybridLogicalClock {
    /// The monotonic system clock, or wall clock.
    time: u128,

    /// The logical clock counter.
    counter: u64,

    /// The last time the monotonic system clock was updated.
    last_physical_time: u64,

    /// The maximum drift allowed for any two clocks in the constellation. This helps enforce linearizability in a very disparate environment.
    max_drift: time::Duration,

    is_locked: AtomicBool,
}

/// Thread-safe implementation of HybridLogicalClock
impl HybridLogicalClock {
    pub fn new(max_drift: time::Duration) -> Self {
        let now = match time::SystemTime::now().duration_since(time::UNIX_EPOCH) {
            Ok(x) => x,
            Err(_) => todo!(),
        }
        .as_nanos();

        // no starting lock.
        let atomic_lock = AtomicBool::new(false);

        HybridLogicalClock {
            time: now,
            counter: 0,
            last_physical_time: 0,
            max_drift: max_drift.clone(),
            is_locked: atomic_lock,
        }
    }

    /// Returns the current time according to the logical clock.
    pub fn now(&mut self) -> (u64, bool) {
        // wait to acquire and then hold lock.
        while self.is_locked.load(Ordering::Acquire) {}
        self.is_locked.store(true, Ordering::Acquire);

        let now = match time::SystemTime::now().duration_since(time::UNIX_EPOCH) {
            Ok(x) => x,
            Err(_) => todo!(),
        }
        .as_nanos();

        if now > self.last_physical_time {
            self.last_physical_time = now;
            self.time = now;
            self.counter += 1;
        } else if now == self.last_physical_time {
            self.counter += 1;
        } else {
            // todo (sienna): don't panic on this.
            self.is_locked.store(false, Release);
            panic!("Time went backwards!");
        }

        self.is_locked.store(false, Release);
        return (self.time + self.counter, true);
    }
}
