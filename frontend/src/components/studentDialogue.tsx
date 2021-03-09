import React, {useState} from "react";
import {Button, Dialog, FormGroup, H2, InputGroup} from "@blueprintjs/core";
import {createStudents, editStudent, TraceStudent} from "../api";
import {onCatch, onCatchPrefix} from "./util";

export function CreateStudentDialogue({isOpen, handleClose}: { isOpen: boolean, handleClose: (studentCreated: boolean) => void }) {
    const [loading, setLoading] = useState(false);

    const onSubmit = (student: TraceStudent) => {
        setLoading(true);
        createStudents([student])
            .then(() => handleClose(true))
            .catch(onCatch)
            .finally(() => setLoading(false))
    }

    return <Dialog
        isOpen={isOpen}
        canEscapeKeyClose={true}
        canOutsideClickClose={true}
        usePortal={true}
        onClose={() => handleClose(false)}
        className="pb-0"
    >
        <StudentEdit
            student={{name: "", email: "", id: "", student_handles: ["", ""]}}
            onSubmit={onSubmit}
            title="Create Student"
            submitButtonText="Create"
            submitButtonLoading={loading}
        />
    </Dialog>;
}

export function EditStudentDialogue({isOpen, handleClose, student}: { isOpen: boolean, student: TraceStudent, handleClose: (studentEdited: boolean) => void }) {
    const [loading, setLoading] = useState(false);

    const onSubmit = (modifiedStudent: TraceStudent) => {
        setLoading(true);
        editStudent(student.id, modifiedStudent)
            .then(() => handleClose(true))
            .catch(onCatch)
            .finally(() => setLoading(false))
    }

    return <Dialog
        isOpen={isOpen}
        canEscapeKeyClose={true}
        canOutsideClickClose={true}
        usePortal={true}
        onClose={() => handleClose(false)}
        className="pb-0"
    >
        <StudentEdit
            student={student}
            onSubmit={onSubmit}
            title={`Edit ${student.name}`}
            submitButtonText="Save"
            submitButtonLoading={loading}
        />
    </Dialog>;
}

export function StudentEdit({student, onSubmit, submitButtonText, title, submitButtonLoading}: {
    student: TraceStudent,
    onSubmit: (student: TraceStudent) => void
    title: string,
    submitButtonText: string
    submitButtonLoading: boolean
}) {
    const [localStudent, setLocalStudent] = useState(student);

    const onSubmitButtonClick = () => {
        // remove empty handles
        onSubmit({...localStudent, student_handles: student.student_handles.filter(handle => handle !== "")});
    }

    return <div className="m-8">
        <H2 className="mb-4">{title}</H2>
        <FormGroup label="Name">
            <InputGroup value={localStudent.name}
                        placeholder="Student's name"
                        leftIcon="person"
                        onChange={(e: any) => setLocalStudent({...localStudent, name: e.target.value})}/>
        </FormGroup>
        <FormGroup label="Email">
            <InputGroup value={localStudent.email}
                        placeholder="Student's email"
                        leftIcon="envelope"
                        onChange={(e: any) => setLocalStudent({...localStudent, email: e.target.value})}/>
        </FormGroup>
        <FormGroup label="Handles">
            {localStudent.student_handles.map((handle, index) => <InputGroup
                    key={index}
                    value={handle}
                    rightElement={<Button icon="trash" onClick={() => {
                        setLocalStudent({
                            ...localStudent,
                            student_handles: localStudent.student_handles.filter((h, idx) => idx !== index)
                        });
                    }}/>}
                    className="mb-1"
                    leftIcon="log-in"
                    placeholder={`Handle ${index + 1}`}
                    onChange={(e: any) => {
                        const newHandles = localStudent.student_handles.map((currHandle, currIndex) => {
                            if (currIndex !== index) {
                                return currHandle;
                            }
                            return e.target.value;
                        });
                        setLocalStudent({...localStudent, student_handles: newHandles});
                    }}
                />
            )}
            <Button icon="add" className="ml-auto" onClick={() => {
                let prevHandles = localStudent.student_handles;
                prevHandles.push('');
                setLocalStudent({...localStudent, student_handles: prevHandles});
            }}>Add</Button>
        </FormGroup>
        <Button intent="primary" className="w-full" onClick={onSubmitButtonClick}
                loading={submitButtonLoading}>{submitButtonText}</Button>
    </div>;
}