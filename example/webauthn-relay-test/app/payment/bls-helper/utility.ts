
import { concatBytes, Hex, hexToNumber, numberToBytesBE, PrivKey } from "@noble/curves/abstract/utils";
import type { ProjPointType, ProjConstructor, AffinePoint } from "@noble/curves/abstract/weierstrass";
import type { Fp2 } from "@noble/curves/abstract/tower";
import { bn254 } from "@noble/curves/bn254"

interface SigPoint {
    px: string,
    py: string,
    pz: string
}

interface PublicPoint {
    px: {
        c0: string,
        c1: string
    },
    py: {
        c0: string,
        c1: string
    }
}

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
        const data: SigPoint = {
            px: hex[0].substring(2),
            py: hex[1].substring(2),
            pz: "01"
        }
        return bn254.G1.ProjectivePoint.fromAffine({
            x: hexToNumber(data.px),
            y: hexToNumber(data.py)
        })
    }
    else {
        return bn254.G1.ProjectivePoint.fromHex("0");
    }
}

function hexToProjPointFp2(hex: string[][]): ProjPointType<Fp2>[] {
    const r: ProjPointType<Fp2>[] = [];
    for (let i = 0; i < hex.length; i++) {
        if (hex[i].length === 4) {
            const data: PublicPoint = {
                px: {
                    c0: hex[i][0].substring(2),
                    c1: hex[i][1].substring(2)
                },
                py: {
                    c0: hex[i][2].substring(2),
                    c1: hex[i][3].substring(2)
                },
            }

            const { Fp2 } = bn254.fields;
            const x = Fp2.fromBigTuple([hexToNumber(data.px.c0), hexToNumber(data.px.c1)]);
            const y = Fp2.fromBigTuple([hexToNumber(data.py.c0), hexToNumber(data.py.c1)]);

            r.push(bn254.G2.ProjectivePoint.fromAffine({
                x, y
            }))
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