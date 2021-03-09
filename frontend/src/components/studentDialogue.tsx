import React, {useState} from "react";
import {Dialog, FormGroup, InputGroup} from "@blueprintjs/core";
import {TraceStudent} from "../api";

export function CreateStudentDialogue({isOpen, handleClose}: { isOpen: boolean, handleClose: () => void }) {
    return <Dialog
        isOpen={isOpen}
        canEscapeKeyClose={true}
        canOutsideClickClose={true}
        usePortal={true}
        onClose={handleClose}
    >
        <StudentEdit student={{name: "", email: "", id: "", student_handles: ["", ""]}} onChange={(student) => {
        }}/>
    </Dialog>;
}

export function StudentEdit({student, onChange}: { student: TraceStudent, onChange: (student: TraceStudent) => void }) {
    const [localStudent, setLocalStudent] = useState(student);

    return <div className="m-8">
        <FormGroup label="Name">
            <InputGroup value={localStudent.name}
                        onChange={(e: any) => setLocalStudent({...localStudent, name: e.target.value})}/>
        </FormGroup>
        <FormGroup label="Email">
            <InputGroup value={localStudent.email}
                        onChange={(e: any) => setLocalStudent({...localStudent, email: e.target.value})}/>
        </FormGroup>
        <FormGroup label="Handles">
            {localStudent.student_handles.map((handle, index) => <InputGroup
                    key={index}
                    value={handle}
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
        </FormGroup>
    </div>;
}