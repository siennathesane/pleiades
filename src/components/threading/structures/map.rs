/*
 *   Pleiades Source Code
 *   Copyright (C) 2023 Sienna Lloyd
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>
 */

use crate::components::threading::structures::traits::Map;
use std::borrow::Borrow;
use std::collections::hash_map;
use std::collections::hash_map::{Entry, Iter, Keys, Values};
use std::hash::Hash;
use std::sync::atomic::Ordering::Relaxed;
use std::sync::atomic::{AtomicBool, Ordering};
use std::thread;

pub struct HashMap<K: Eq + PartialEq + Hash, V> {
    data: hash_map::HashMap<K, V>,
    locked: AtomicBool,
}

impl<K: Eq + PartialEq + Hash, V> Map<K, V> for HashMap<K, V> {
    fn new() -> Self {
        Self {
            data: hash_map::HashMap::new(),
            locked: AtomicBool::new(false),
        }
    }

    fn with_capacity(capacity: usize) -> Self {
        Self {
            data: hash_map::HashMap::with_capacity(capacity),
            locked: AtomicBool::new(false),
        }
    }

    fn contains_key<Q>(&self, k: &Q) -> bool
    where
        K: Borrow<Q>,
        Q: Hash + Eq + ?Sized,
    {
        self.data.contains_key(k)
    }

    fn entry(&mut self, key: K) -> Entry<'_, K, V> {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        let val = self.data.entry(key);
        self.locked.store(false, Relaxed);

        val
    }

    fn get<Q>(&self, key: &Q) -> Option<&V>
    where
        K: Borrow<Q>,
        Q: Hash + Eq + ?Sized,
    {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        let val = self.data.get(key);
        self.locked.store(false, Relaxed);

        val
    }

    fn insert(&mut self, k: K, v: V) -> Option<V> {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        let val = self.data.insert(k, v);
        self.locked.store(false, Relaxed);

        val
    }

    fn iter(&self) -> Iter<'_, K, V> {
        return self.data.iter().clone();
    }

    fn remove<Q>(&mut self, k: &Q) -> Option<V>
    where
        K: Borrow<Q>,
        Q: Hash + Eq + ?Sized,
    {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        let val = self.data.remove(k);
        self.locked.store(false, Relaxed);

        val
    }

    fn keys(&self) -> Keys<'_, K, V> {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);
        let val = self.data.keys().clone();

        self.locked.store(false, Relaxed);

        val
    }

    fn values(&self) -> Values<'_, K, V> {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);
        let val = self.data.values().clone();

        self.locked.store(false, Relaxed);

        val
    }
}
