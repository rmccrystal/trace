import React, {Dispatch, SetStateAction, useEffect, useState} from 'react';
import './app.scss';
import Scan from './components/scan';
import Nav from "./components/nav";
import {createGlobalState} from "react-hooks-global-state";
import * as Api from "./api";
import {Card, FocusStyleManager, Spinner} from "@blueprintjs/core";
import {Switch, Route, Redirect, BrowserRouter} from "react-router-dom";
import Dashboard from "./components/dashboard";
import {getLocations, TraceLocation} from "./api";
import useLocalStorage, {onCatchPrefix} from "./components/util";

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

    let [dark, setDark] = useLocalStorage('dark', false);
    let [_location, setLocation] = useState<TraceLocation | null>(null);
    let [locationID, setLocationID] = useLocalStorage('location', '');

    useEffect(() => {
        if (_location) {
            setLocationID(_location.id)
        }
    }, [_location, setLocationID])

    useEffect(() => {
        getLocations()
            .then(locations => {
                setLocation(locations.find(value => value.id === locationID) || locations[0])
            })
            .catch(onCatchPrefix("Error getting list of locations: "));
    }, [locationID])

    if (_location === null) {
        return <Spinner className="m-auto absolute inset-0"/>
    }
    let location = _location!;

    return (
        <BrowserRouter>
            <div className={`app flex flex-col items-center content-center ${dark ? "bp3-dark" : "bp3-light"}`}>
                <Nav setLocation={setLocation} location={location} onToggleDark={() => setDark(!dark)}/>
                <Switch>
                    <Route exact path="/scan" component={() => <Scan location={location} className="p-16 m-8 my-auto" elevation={0}/>}/>
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