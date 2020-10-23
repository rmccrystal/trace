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
    location: TraceLocation,
    student: TraceStudent,
    time: Date,
    event_type: EventType
}

export async function scan(student_handle: string, location_id: string): Promise<TraceEvent> {
    return await sendApiRequest<TraceEvent>("POST", "scan", {student_handle, location_id});
}

export interface TraceLocation {
    id: string,
    name: string,
    timeout: number
}

export async function getLocations(): Promise<TraceLocation[]> {
    return await sendApiRequest<TraceLocation[]>("GET", "location");
}

export interface TraceStudent {
    id: string,
    name: string,
    email: string,
    student_handles: string[]
}

export async function getStudentsAtLocation(location_id: string): Promise<{student: TraceStudent, time: Date}[]> {
    let data = await sendApiRequest<{student: TraceStudent, time: Date}[]>("GET", `location/${location_id}/students`);
    data.map((st) => {
        st.time = new Date(st.time);
    });
    return data;
}

export async function logoutStudent(student_id: string, location_id: string): Promise<TraceEvent> {
    return await sendApiRequest("POST", `student/${student_id}/logout`, {location_id: location_id});
}

export async function getStudents(): Promise<TraceStudent[]> {
    return await sendApiRequest("GET", `student`);
}

export async function logoutAll(location_id: string): Promise<null> {
    return await sendApiRequest("POST", `location/${location_id}/logoutAll`);
}

export async function createNewLocation(location: Partial<TraceLocation>): Promise<TraceLocation> {
    return await sendApiRequest("POST", `location`, location);
}

export async function createStudents(students: TraceStudent[]): Promise<TraceStudent[]> {
    return await sendApiRequest("POST", "students", students);
}

export interface ContactReport {
    target_student: TraceStudent,
    contacts: {
        student: TraceStudent,
        duration: number, // duration in ms
        depth: number
    }[]
}
// duration should be in ms
export async function getStudentContacts(student_id: string, duration: number, max_depth: number = 1): Promise<ContactReport> {
    return await sendApiRequest("GET", `student/${student_id}/contacts`, {duration, max_depth});
}

export interface LocationVisit {
    student: TraceStudent,
    leave_time: Date,
    enter_time: Date,
}
export async function getLocationVisits(location_id: string): Promise<LocationVisit[]> {
    let data = await sendApiRequest<LocationVisit[]>("GET", `location/${location_id}/visits`)
    data.map(el => {
        el.leave_time = new Date(el.leave_time);
        el.enter_time = new Date(el.enter_time);
        return el
    })
    return data;
}