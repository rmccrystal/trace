import React, {useEffect, useState} from "react";
import {Card, HTMLTable, IHTMLTableProps} from "@blueprintjs/core";
import {getStudents, Student} from "../api";
import {onCatch} from "./util";

export default function StudentList() {
    const [students, setStudents] = useState<Student[]>([]);

    useEffect(() => {
        getStudents()
            .then(setStudents)
            .catch(onCatch)
    }, []);

    return <Card className="mt-8 max-w-3xl w-full p-0">
        <StudentTable students={students} className="w-full" striped bordered/>
    </Card>
}

function StudentTable({students, ...props}: { students: Student[] } & IHTMLTableProps) {
    return <HTMLTable {...props}>
        <thead>
        <tr>
            <th>Name</th>
            <th>Email</th>
            <th>Handles</th>
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