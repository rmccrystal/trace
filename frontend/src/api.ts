import axios from 'axios';

export enum EventType {
    Enter,
    Leave
}

export interface TraceEvent {
    id: string,
    location_id: string,
    location_name: string,
    student_id: string,
    student_name: string,
    time: Date,
    event_type: EventType
}

// sendApiRequest sends a request to the origin with the method and relative path to the api.
// For example, if path was scan this function will request /api/scan
// It will return a promise with the type of the response. If the response contains
// an error, the error will be thrown
async function sendApiRequest<T>(method: "GET" | "PUT" | "POST" | "DELETE" | "PATCH", path: string, data?: any): Promise<T> {
    path = `/api/${path}`;
    const resp = await axios.request({method, baseURL: path, data, validateStatus: () => true});
    const success: boolean = resp.data.success;
    if (success) {
        return resp.data.data as T;
    } else {
        if (resp.data.error) {
            throw new Error(resp.data.error);
        } else {
            throw new Error(`Received error code from server: ${resp.status}`);
        }
    }
}

export async function scan(student_handle: string, location_id: string): Promise<TraceEvent> {
    return await sendApiRequest<TraceEvent>("POST", "scan", {student_handle, location_id});
}