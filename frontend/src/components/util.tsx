import {Toaster} from "@blueprintjs/core";
import {useState} from "react";

const errorToaster = Toaster.create({position: "top", maxToasts: 4});

// onCatch should be called whenever an exception during
// a request occurs using .catch on the promise
export function onCatch(reason: any) {
    errorToaster.show({message: reason.toString(), icon: "error", intent: "danger"})
}

// Returns a function similar to onCatch with a prefix
export function onCatchPrefix(prefix: string): (reason: any) => void {
    return reason => {
        errorToaster.show({message: prefix + reason.toString(), icon: "error", intent: "danger"})
    }
}

export function formatAMPM(date: Date): string {
    var hours = date.getHours();
    var minutes: string | number = date.getMinutes();
    var ampm = hours >= 12 ? 'PM' : 'AM';
    hours = hours % 12;
    hours = hours ? hours : 12; // the hour '0' should be '12'
    minutes = minutes < 10 ? '0'+minutes : minutes;
    var strTime = hours + ':' + minutes + ' ' + ampm;
    return strTime;
}

// From https://usehooks-typescript.com/use-local-storage/
export default function useLocalStorage<T>(key: string, initialValue: T) {
    // Prevent build error "window is undefined" but keep keep working
    const isServer = typeof window === 'undefined'
    // State to store our value
    // Pass initial state function to useState so logic is only executed once
    const [storedValue, setStoredValue] = useState(() => {
        // Get from local storage then
        // parse stored json or return initialValue
        if (isServer) {
            return initialValue
        }
        try {
            const item = window.localStorage.getItem(key)
            return item ? JSON.parse(item) : initialValue
        } catch (error) {
            console.log(error)
            return initialValue
        }
    })
    // Return a wrapped version of useState's setter function that ...
    // ... persists the new value to localStorage.
    const setValue = (value: T) => {
        try {
            // Allow value to be a function so we have same API as useState
            const valueToStore =
                value instanceof Function ? value(storedValue) : value
            // Save state
            setStoredValue(valueToStore)
            // Save to local storage
            if (!isServer) {
                window.localStorage.setItem(key, JSON.stringify(valueToStore))
            }
        } catch (error) {
            console.log(error)
        }
    }
    return [storedValue, setValue]
}
