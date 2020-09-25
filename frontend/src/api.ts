import axios from 'axios';


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

export async function scan(student_handle: string, location_id: string): Promise<TraceEvent> {
    return await sendApiRequest<TraceEvent>("POST", "scan", {student_handle, location_id});
}

export interface Location {
    id: string,
    name: string,
    timeout: number
}

export async function getLocations(): Promise<Location[]> {
    return await sendApiRequest<Location[]>("GET", "location");
}

export interface Student {
    id: string,
    name: string,
    email: string,
    student_handles: string[]
}

export async function getStudentsAtLocation(location_id: string): Promise<{student: Student, event: TraceEvent}[]> {
    let data = await sendApiRequest<{student: Student, event: TraceEvent}[]>("GET", `location/${location_id}/students`);
    data.map((st) => {
        st.event.time = new Date(st.event.time);
    });
    return data;
}

export async function logoutStudent(student_id: string, location_id: string): Promise<TraceEvent> {
    return await sendApiRequest("POST", `student/${student_id}/logout`, {location_id: location_id});
}