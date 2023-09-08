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

use crate::components::threading::structures::traits::List;
use std::ops::Index;
use std::slice::SliceIndex;
use std::sync::atomic::{AtomicBool, Ordering};
use std::thread;

pub struct Vector<T> {
    data: Vec<T>,
    locked: AtomicBool,
}

impl<T> List<T> for Vector<T> {
    fn new() -> Self {
        Vector {
            data: Vec::new(),
            locked: AtomicBool::new(false),
        }
    }

    fn with_capacity(capacity: usize) -> Self {
        Vector {
            data: Vec::with_capacity(capacity),
            locked: AtomicBool::new(false),
        }
    }

    fn push(&mut self, value: T) {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        self.data.push(value);
        self.locked.store(false, Ordering::Relaxed);
    }

    fn pop(&mut self) -> T {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        // fetch the value, then unlock
        let val = self.data.pop().unwrap();
        self.locked.store(false, Ordering::Relaxed);

        val
    }

    fn len(&self) -> usize {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        let val = self.data.len();
        self.locked.store(false, Ordering::Relaxed);

        val
    }

    fn capacity(&self) -> usize {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        let val = self.data.capacity();
        self.locked.store(false, Ordering::Relaxed);

        val
    }

    fn is_empty(&self) -> bool {
        self.data.is_empty()
    }

    fn clear(&mut self) {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        self.data.clear();

        self.locked.store(false, Ordering::Relaxed);
    }

    fn remove(&mut self, index: usize) -> T {
        // yield the thread until we can acquire the lock
        while self.locked.load(Ordering::Acquire) {
            thread::yield_now();
        }
        self.locked.store(true, Ordering::Relaxed);

        let val = self.data.remove(index);

        self.locked.store(false, Ordering::Relaxed);

        val
    }
}

impl<T, Idx> Index<Idx> for Vector<T>
where
    Idx: SliceIndex<[T], Output = T>,
{
    type Output = T;

    fn index(&self, index: Idx) -> &Self::Output {
        &self.data[index]
    }
}

#[cfg(test)]
mod tests {
    use crate::components::threading::structures::list::Vector;
    use crate::components::threading::structures::traits::List;

    use rand::Rng;
    use std::sync::{Arc, Mutex};
    use std::thread;
    use std::thread::JoinHandle;
    use std::time::Duration;

    #[test]
    fn new() {
        let mut list: Vector<i32> = Vector::new();
        list.push(1);
        assert_eq!(list.len(), 1);
    }

    #[test]
    fn with_capacity() {
        let list: Vector<u128> = Vector::with_capacity(10);
        assert_eq!(list.capacity(), 10);
    }

    #[test]
    fn push() {
        let mut _thread_list: Vector<JoinHandle<()>> = Vector::new();
        let shared_list: Arc<Mutex<Vector<u64>>> = Arc::new(Mutex::new(Vector::new()));

        let threads = 4;
        let additions = 100;

        for _i in 0..threads {
            let my_list = shared_list.clone();
            let handle = thread::spawn(move || {
                let mut shared = my_list.lock().unwrap();

                // generate some random jitter for the sleep function so we can eliminate data races
                let sleep_val = rand::thread_rng().gen_range(0..10);
                for i in 0..additions {
                    shared.push(i);
                    thread::sleep(Duration::from_micros(sleep_val * i));
                }
            });
            _thread_list.push(handle);
        }

        for handle in _thread_list.data {
            let _ = handle.join();
        }

        assert_eq!(
            shared_list.clone().lock().unwrap().len(),
            threads * additions as usize
        );
    }

    #[test]
    fn pop() {
        let mut _thread_list: Vector<JoinHandle<()>> = Vector::new();
        let shared_list: Arc<Mutex<Vector<u64>>> = Arc::new(Mutex::new(Vector::new()));

        let threads = 4;
        let additions = 100;

        for _i in 0..threads {
            let my_list = shared_list.clone();
            let handle = thread::spawn(move || {
                let mut shared = my_list.lock().unwrap();

                // generate some random jitter for the sleep function so we can eliminate data races
                let sleep_val = rand::thread_rng().gen_range(0..10);
                for i in 0..additions {
                    shared.push(i);
                    thread::sleep(Duration::from_micros(sleep_val * i));
                    shared.pop();
                }
            });
            _thread_list.push(handle);
        }

        for handle in _thread_list.data {
            let _ = handle.join();
        }

        assert_eq!(shared_list.clone().lock().unwrap().len(), 0);
    }

    #[test]
    fn len() {
        let mut list: Vector<i32> = Vector::new();
        list.push(1);
        assert_eq!(list.len(), 1);
    }

    #[test]
    fn is_empty() {
        let mut list: Vector<i32> = Vector::new();
        list.push(1);
        assert_eq!(list.is_empty(), false);

        list.pop();
        assert_eq!(list.is_empty(), true);
    }

    #[test]
    fn clear() {
        let mut list: Vector<i32> = Vector::new();

        for i in 0..99 {
            list.push(i);
        }

        list.clear();
        assert_eq!(list.len(), 0);
    }

    #[test]
    fn remove() {
        let mut list: Vector<i32> = Vector::new();

        for i in 0..99 {
            list.push(i);
        }

        for i in 0..49 {
            list.remove(i);
        }

        assert_eq!(list.len(), 50);
    }
}
