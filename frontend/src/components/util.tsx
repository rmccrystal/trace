import {Toaster} from "@blueprintjs/core";

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
