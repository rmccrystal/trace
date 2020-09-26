import React, {Dispatch, SetStateAction, useEffect, useState} from 'react';
import './app.scss';
import Scan from './components/scan';
import Nav from "./components/nav";
import {createGlobalState} from "react-hooks-global-state";
import * as Api from "./api";
import {FocusStyleManager, Spinner} from "@blueprintjs/core";
import {Switch, Route, Redirect, BrowserRouter} from "react-router-dom";
import Dashboard from "./components/dashboard";
import {getLocations, TraceLocation} from "./api";
import {onCatchPrefix} from "./components/util";

export interface GlobalState {
    dark: boolean
}

export const {useGlobalState} = createGlobalState<GlobalState>({
    dark: false
});

function App() {
    useEffect(() => {
        FocusStyleManager.onlyShowFocusOnTabs();
    }, [])

    let [dark] = useGlobalState('dark');
    let [_location, setLocation] = useState<TraceLocation | null>(null);

    useEffect(() => {
        getLocations()
            .then(locations => setLocation(locations[0]))
            .catch(onCatchPrefix("Error getting list of locations: "));
    }, [])

    if (_location === null) {
        return <Spinner className="m-auto absolute inset-0" />
    }
    let location = _location!;

    return (
        <BrowserRouter>
            <div className={`app flex flex-col items-center content-center ${dark ? "bp3-dark" : "bp3-light"}`}>
                <Nav setLocation={setLocation} location={location}/>
                <Switch>
                    <Route exact path="/scan" component={() => <Scan location={location}/>}/>
                    <Route path="/dashboard" component={() => <Dashboard location={location}/>}/>
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