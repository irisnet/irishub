#pragma once
#ifndef __BIN2ASCII_H__
#define __BIN2ASCII_H__

#include <string>
#include <stdexcept>
    
inline std::string Hex2Bin(const std::string &sValue)
{
    if (sValue.size() % 2)
        throw std::runtime_error("Odd hex data size");
    static const char sLookUp[] = ""
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x00
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x10
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x20
        "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x80\x80\x80\x80\x80\x80" // 0x30
        "\x80\x0a\x0b\x0c\x0d\x0e\x0f\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x40
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x50
        "\x80\x0a\x0b\x0c\x0d\x0e\x0f\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x60
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x70
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x80
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x90
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xa0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xb0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xc0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xd0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xe0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xf0
        "";
    std::string sRet;
    sRet.reserve(sValue.size() / 2);
    for (size_t i = 0; i < sValue.size(); i += 2) {
        char cHigh = sLookUp[sValue[i]];
        char cLow = sLookUp[sValue[i+1]];
        if (0x80 & (cHigh | cLow))
            throw std::runtime_error("Invalid hex data: " + sValue.substr(i, 6));
        sRet.push_back((cHigh << 4) | cLow);
    }
    return sRet;
}

inline std::string Bin2Hex(const std::string &sValue)
{
    static const char sLookUp[] = "0123456789abcdef";
    std::string sRet;
    sRet.reserve(sValue.size() * 2);
    for (size_t i = 0; i < sValue.size(); i++) {
        char cHigh = sValue[i] >> 4;
        char cLow = sValue[i] & 0xf;
        sRet.push_back(sLookUp[cHigh]);
        sRet.push_back(sLookUp[cLow]);
    }
    return sRet;
}

inline std::string B64Encode(const std::string &sValue)
{
    typedef unsigned char uc;
    static const char sLookUp[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    const uc * sData = (const uc *) sValue.c_str();
    std::string sRet;
    sRet.reserve(sValue.size() * 4 / 3 + 3);
    for (size_t i = 0; i < sValue.size(); i += 3) {
        unsigned n = sData[i] << 16;
        if (i + 1 < sValue.size()) n |= sData[i + 1] << 8;
        if (i + 2 < sValue.size()) n |= sData[i + 2];

        uc n0 = (uc)(n >> 18) & 0x3f;
        uc n1 = (uc)(n >> 12) & 0x3f;
        uc n2 = (uc)(n >>  6) & 0x3f;
        uc n3 = (uc)(n      ) & 0x3f;

        sRet.push_back(sLookUp[n0]);
        sRet.push_back(sLookUp[n1]);
        if (i + 1 < sValue.size()) sRet.push_back(sLookUp[n2]);
        if (i + 2 < sValue.size()) sRet.push_back(sLookUp[n3]);
    }
    for (int i = 0; i < (3 - sValue.size() % 3) % 3; i++)
        sRet.push_back('=');
    return sRet;
}

inline std::string B64Decode(const std::string &sValue)
{
    typedef unsigned char uc;
    static const char sLookUp[] = ""
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x00
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x10
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x3e\x80\x80\x80\x3f" // 0x20
        "\x34\x35\x36\x37\x38\x39\x3a\x3b\x3c\x3d\x80\x80\x80\x00\x80\x80" // 0x30
        "\x80\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e" // 0x40
        "\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x80\x80\x80\x80\x80" // 0x50
        "\x80\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28" // 0x60
        "\x29\x2a\x2b\x2c\x2d\x2e\x2f\x30\x31\x32\x33\x80\x80\x80\x80\x80" // 0x70
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x80
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0x90
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xa0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xb0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xc0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xd0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xe0
        "\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" // 0xf0
        "";
    std::string sRet;
    if (!sValue.size()) return sRet;
    if (sRet.size() % 4)
        throw std::runtime_error("Invalid base64 data size");
    size_t iPad = 0;
    if (sValue[sValue.size() - 1] == '=') iPad++;
    if (sValue[sValue.size() - 2] == '=') iPad++;

    sRet.reserve(sValue.size() * 3 / 4 + 3);
    for (size_t i = 0; i < sValue.size(); i += 4) {
        uc n0 = sLookUp[(uc) sValue[i+0]];
        uc n1 = sLookUp[(uc) sValue[i+1]];
        uc n2 = sLookUp[(uc) sValue[i+2]];
        uc n3 = sLookUp[(uc) sValue[i+3]];
        if (0x80 & (n0 | n1 | n2 | n3))
            throw std::runtime_error("Invalid hex data: " + sValue.substr(i, 4));
        unsigned n = (n0 << 18) | (n1 << 12) | (n2 << 6) | n3;
        sRet.push_back((n >> 16) & 0xff);
        if (sValue[i+2] != '=') sRet.push_back((n >> 8) & 0xff);
        if (sValue[i+3] != '=') sRet.push_back((n     ) & 0xff);
    }
    return sRet;
}

#endif//__BIN2ASCII_H__
