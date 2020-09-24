import React, {Dispatch, SetStateAction, useEffect, useState} from 'react';
import './app.scss';
import Scan from './components/scan';
import Nav from "./components/nav";
import {createGlobalState} from "react-hooks-global-state";
import * as Api from "./api";
import {FocusStyleManager} from "@blueprintjs/core";
import {Switch, Route, Redirect} from "react-router-dom";
import StudentList from "./components/studentList";

export interface GlobalState {
    location?: Api.Location
    dark: boolean
}

export const {useGlobalState} = createGlobalState<GlobalState>({
    location: undefined, dark: false
});

function App() {
    useEffect(() => {
        FocusStyleManager.onlyShowFocusOnTabs();
    }, [])

    let [dark] = useGlobalState('dark');

    return (
        <div className={`app ${dark ? "bp3-dark" : "bp3-light"}`}>
            <Nav/>
            <Switch>
                <Route exact path="/" component={Scan} />
                <Route path="/students" component={StudentList} />
                <Redirect from="*" to="/" />
            </Switch>
        </div>
    );
}

export default App;