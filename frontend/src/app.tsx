import React, {Dispatch, SetStateAction, useState} from 'react';
import './app.scss';
import Scan from './components/scan';
import Nav from "./components/nav";
import {createGlobalState} from "react-hooks-global-state";
import * as Api from "./api";

export interface GlobalState {
    location?: Api.Location
}

export const { useGlobalState } = createGlobalState<GlobalState>({location: undefined});

function App() {
    return (
        <div className="app">
            <Nav/>
            <Scan/>
        </div>
    );
}

export default App;