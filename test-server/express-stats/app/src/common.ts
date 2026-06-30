import { Request, Response, NextFunction } from 'express'


export function calculateDelay(delayMin: number, deplayMax: number | undefined): number {

    if (deplayMax) {
        const rangeFrom = Math.min(delayMin, deplayMax);
        const rangeTo = Math.max(delayMin, deplayMax);

        return Math.round(Math.random() * (rangeTo - rangeFrom) + rangeFrom);
    }

    return delayMin;
}

export function randomCode(): number {
    const codeHi = Math.round(Math.random() * (5 - 1) + 1);
    const codeLo = Math.round(Math.random() * 4);

    return codeHi * 100 + codeLo;
}