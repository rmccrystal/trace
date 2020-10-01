import React, {useEffect, useState} from "react";
import {
    Button,
    Card, Dialog,
    FileInput,
    HTMLTable,
    ICardProps,
    Icon,
    IHTMLTableProps,
    Spinner,
    Tag,
    Tooltip
} from "@blueprintjs/core";
import {getStudents, Student} from "../api";
import {onCatch} from "./util";
import {CSVReader} from "react-papaparse";

export default function ManageStudents({...props}: ICardProps) {
    const [students, setStudents] = useState<Student[]>([]);
    const [loading, setLoading] = useState(true);
    const [uploadDialogOpen, setUploadDialogOpen] = useState(false);

    useEffect(() => {
        getStudents()
            .then(students => {
                setStudents(students);
                setLoading(false);
            })
            .catch(onCatch)
    }, []);

    return <Card className="m-8 max-w-3xl w-full" {...props}>
        <div className="flex flex-row gap-2">
            <Button className="w-full mb-4" onClick={() => setUploadDialogOpen(true)}>Add students from CSV</Button>
        </div>
        <UploadStudentCsvDialog isOpen={uploadDialogOpen} onSubmit={locations => {
            console.log(locations.toString())
        }} onClose={() => setUploadDialogOpen(false)}/>
        {
            loading
                ? <Spinner className="m-8"/>
                : <Card className="p-0 w-full" elevation={1}>
                    <StudentTable students={students} className="w-full" striped bordered/>
                </Card>
        }
    </Card>
}

function UploadStudentCsvDialog({isOpen, onSubmit, onClose}: { isOpen: boolean, onClose: () => void, onSubmit: (students: Student[]) => void }) {
    const handleFileLoad = (data: any) => {
        console.log(data);
    }

    return <Dialog isOpen={isOpen} title="Upload students from CSV" onClose={onClose}>
        <span className="m-8">
            <CSVReader onFileLoad={handleFileLoad}>
                <span>Drag CSV file here</span>
            </CSVReader>
        </span>
    </Dialog>
}

export function StudentTable({students, loading, ...props}: { students: Student[], loading?: boolean } & IHTMLTableProps) {
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

function StudentRow({student}: { student: Student }) {
    return <tr>
        <td>{student.name || "-"}</td>
        <td>{student.email || "-"}</td>
        <td>{student.student_handles.join(", ")}</td>
    </tr>
}