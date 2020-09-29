import React, {useEffect, useState} from "react";
import {Card, HTMLTable, Icon, IHTMLTableProps, Spinner, Tag, Tooltip} from "@blueprintjs/core";
import {getStudents, Student} from "../api";
import {onCatch} from "./util";

export default function StudentList() {
    const [students, setStudents] = useState<Student[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        getStudents()
            .then(students => {
                setStudents(students);
                setLoading(false);
            })
            .catch(onCatch)
    }, []);

    return <Card className="mt-8 max-w-3xl w-full p-0">
        <StudentTable students={students} className="w-full" striped bordered loading={loading}/>
    </Card>
}

function StudentTable({students, loading, ...props}: { students: Student[], loading: boolean } & IHTMLTableProps) {
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