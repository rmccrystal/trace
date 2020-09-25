import React, {Dispatch, SetStateAction, useEffect, useState} from 'react';
import './app.scss';
import Scan from './components/scan';
import Nav from "./components/nav";
import {createGlobalState} from "react-hooks-global-state";
import * as Api from "./api";
import {FocusStyleManager} from "@blueprintjs/core";
import {Switch, Route, Redirect, BrowserRouter} from "react-router-dom";
import Dashboard from "./components/dashboard";

export interface GlobalState {
    location?: Api.TraceLocation
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
        <BrowserRouter>
            <div className={`app flex flex-col items-center content-center ${dark ? "bp3-dark" : "bp3-light"}`}>
                <Nav/>
                <Switch>
                    <Route exact path="/scan" component={Scan}/>
                    <Route path="/dashboard" component={Dashboard}/>
                    <Redirect from="*" to="/scan"/>
                </Switch>
                <div className="mt-auto text-center mb-3 bp3-text-muted">
                    Created by Ryan McCrystal | <a href="https://github.com/rmccrystal/trace" target="_blank">github</a>
                </div>
            </div>
        </BrowserRouter>
    );
}

export default App;