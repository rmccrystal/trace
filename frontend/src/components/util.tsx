import {Toaster} from "@blueprintjs/core";

const errorToaster = Toaster.create({position: "top", maxToasts: 4});

// onCatch should be called whenever an exception during
// a request occurs using .catch on the promise
export function onCatch(reason: any) {
    errorToaster.show({message: reason.toString(), icon: "error", intent: "danger"})
}