import React, {Dispatch, SetStateAction, useEffect, useState} from 'react';
import './app.scss';
import Scan from './components/scan';
import Nav from "./components/nav";
import {createGlobalState} from "react-hooks-global-state";
import * as Api from "./api";
import {FocusStyleManager} from "@blueprintjs/core";

export interface GlobalState {
    location?: Api.Location
}

export const { useGlobalState } = createGlobalState<GlobalState>({location: undefined});

function App() {
    useEffect(() => {
        FocusStyleManager.onlyShowFocusOnTabs();
    }, [])
    return (
        <div className="app">
            <Nav/>
            <Scan/>
        </div>
    );
}

export default App;