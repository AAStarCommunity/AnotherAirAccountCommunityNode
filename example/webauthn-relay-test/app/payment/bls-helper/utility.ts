
import { concatBytes, Hex, numberToBytesBE, PrivKey } from "@noble/curves/abstract/utils";
import { bn254 } from '@kevincharm/noble-bn254-drand'
import type { ProjPointType, ProjConstructor, AffinePoint } from "@noble/curves/abstract/weierstrass";
import type { Fp2 } from "@noble/curves/abstract/tower";

const getBigIntPoint = (point: ProjPointType<bigint>) => {
    return concatBytes(
        numberToBytesBE(point.x, 32),
        numberToBytesBE(point.y, 32),
    )
}

const getFp2Point = (point: ProjPointType<Fp2>) => {
    return concatBytes(
        numberToBytesBE(point.x.c1, 32),
        numberToBytesBE(point.x.c0, 32),
        numberToBytesBE(point.y.c1, 32),
        numberToBytesBE(point.y.c0, 32),
    )
}

function hexToProjPoint(hex: string[]): ProjPointType<bigint> {
    if (hex.length === 2) {
        return bn254.G1.ProjectivePoint.fromHex(hex[0].substring(2) + hex[1].substring(2));
    }
    else {
        return bn254.G1.ProjectivePoint.fromHex("0");
    }
}

function hexToProjPointFp2(hex: string[][]): ProjPointType<Fp2>[] {
    const r: ProjPointType<Fp2>[] = [];
    for (let i = 0; i < hex.length; i++) {
        if (hex[i].length === 4) {
            const k1 = hex[i][0].substring(2) + hex[i][1].substring(2)
            const k2 = hex[i][2].substring(2) + hex[i][3].substring(2);
            const b1 = bn254.G2.ProjectivePoint.fromHex(k1);
            console.log(b1);
            r.push(bn254.G2.ProjectivePoint.fromHex(k1));
        }
        else {
            r.push(bn254.G2.ProjectivePoint.fromHex("0"));
            break
        }
    }
    return r
}


export const getAggr = (signs: string[], publicKeys: string[][], messages: string[]) => {
    const aggSignature = hexToProjPoint(signs);
    const hm = hexToProjPoint(messages);
    const publicPoints = hexToProjPointFp2(publicKeys);

    let calldata = concatBytes(getBigIntPoint(aggSignature), getFp2Point(bn254.G2.ProjectivePoint.BASE));
    for (let i = 0; i < publicPoints.length; i++) {
        calldata = concatBytes(calldata, getBigIntPoint(hm), getFp2Point(publicPoints[i].negate()));
    }
    return calldata
};