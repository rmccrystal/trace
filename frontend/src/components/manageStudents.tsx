import React, {useEffect, useState} from "react";
import {Card, FileInput, HTMLTable, ICardProps, Icon, IHTMLTableProps, Spinner, Tooltip} from "@blueprintjs/core";
import {createStudents, getStudents, TraceStudent} from "../api";
import {onCatch} from "./util";
import Papa from "papaparse";

export default function ManageStudents({...props}: ICardProps) {
    const [students, setStudents] = useState<TraceStudent[]>([]);
    const [loading, setLoading] = useState(true);

    const updateStudents = () => {
        getStudents()
            .then(students => {
                setStudents(students);
                setLoading(false);
            })
            .catch(onCatch)
    }

    useEffect(() => {
        updateStudents();
    }, []);

    const onStudentUpload = (students: TraceStudent[]) => {
        createStudents(students)
            .then(() => {
                updateStudents();
            })
            .catch(onCatch);
    }

    return <Card className="m-8 max-w-3xl w-full" {...props}>
        <div className="flex flex-row gap-2 mb-4">
            <StudentCSVUpload onUpload={onStudentUpload} />
        </div>
        {
            loading
                ? <Spinner className="m-8"/>
                : <Card className="p-0 w-full" elevation={1}>
                    <StudentTable students={students} className="w-full" condensed striped bordered/>
                </Card>
        }
    </Card>
}

function StudentCSVUpload({onUpload}: {onUpload: (students: TraceStudent[]) => void}) {
    const handleError = onCatch

    const handleStudentUploadChange: React.FormEventHandler<HTMLInputElement> = event => {
        const file = (event.target as any).files[0]
        if (!file) {
            return
        }

        const reader = new FileReader();
        reader.readAsText(file);

        reader.onerror = handleError;
        reader.onload = ev => {
            if (!ev.target) {
                handleError("Error reading CSV file");
                return
            }

            const text = ev.target.result?.toString();
            if (!text) {
                handleError("No file content found");
                return
            }
            if (!text.split('\n')[0].startsWith("name,email,handle")) {
                handleError("CSV header must be name,email,handle");
                return
            }

            // We can't pass arrays as CSV so we will have to use a single handle
            interface CSVStudent extends TraceStudent {
                handle: string
            }

            // Parse csv into students
            const results = Papa.parse<CSVStudent>(text, {
                header: true,
            });

            if (results.errors.length > 0) {
                handleError(`Encountered one or more errors parsing the CSV file: ${
                    results.errors.map(e => e.message).join(', ')
                }`);
                return
            }

            const csvStudents = results.data;
            if (csvStudents.length === 0) {
                handleError("No students were found in the CSV")
                return
            }

            const students: TraceStudent[] = csvStudents.map(st => {
                st.student_handles = [st.handle];
                return st as TraceStudent;
            })

            onUpload(students);
        }
    }

    return <FileInput
        text="Upload students from CSV (header should be name,email,handle)"
        className="w-full"
        inputProps={{accept: ".csv"}}
        onInputChange={handleStudentUploadChange}
    />
}

export function StudentTable({students, loading, ...props}: { students: TraceStudent[], loading?: boolean } & IHTMLTableProps) {
    if (loading) {
        return <Spinner className="m-8"/>
    }

    return <HTMLTable {...props}>
        <thead>
        <tr>
            <th>Name</th>
            <th>Email</th>
            <th>
                <Tooltip content="A student handle is text that can be entered into the scan tab to log in or log out">
                    <span>Handles <Icon style={{verticalAlign: "top"}} icon="help" iconSize={10}/></span>
                </Tooltip>
            </th>
        </tr>
        </thead>
        <tbody>
        {students.map(student => <StudentRow student={student} key={student.id}/>)}
        </tbody>
    </HTMLTable>
}

function StudentRow({student}: { student: TraceStudent }) {
    return <tr>
        <td>{student.name || "-"}</td>
        <td>{student.email || "-"}</td>
        <td>{student.student_handles?.join(", ") || "-"}</td>
    </tr>
}