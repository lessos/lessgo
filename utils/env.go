package utils

import (
    "errors"
    "os/exec"
    "runtime"
    "strings"
)

const (
    DistGen = "gen"
    ArchAll = "noarch"
)

type EnvOs struct {
    Kernel string
}

func EnvDistArch() (string, string, error) {

    cmd, err := exec.LookPath("lsb_release")
    if err != nil {
        return "", "", err
    }

    rs, err := exec.Command(cmd, "-r", "-i", "-s").Output()
    if err != nil {
        return "", "", err
    }

    dist := ""
    arch := runtime.GOARCH

    out := strings.Replace(string(rs), "\n", " ", -1)
    rs2 := strings.Split(out, " ")
    if len(rs2) < 2 {
        return dist, arch, errors.New("Unknow ENV")
    }

    if rs2[0] == "CentOS" {
        dist = "el"
    } else if rs2[0] == "Debian" {
        dist = "de"
    } else {
        return dist, arch, errors.New("Unknow ENV")
    }

    ver := strings.Split(rs2[1], ".")
    if len(ver) == 0 {
        return dist, arch, errors.New("Unknow ENV")
    }
    if ver[0] == "6" {
        dist += "6"
    } else {
        return dist, arch, errors.New("Unknow ENV")
    }

    if arch != "amd64" {
        return dist, arch, errors.New("Unknow ENV")
    }

    return dist, arch, nil
}

func EnvOsInfo() EnvOs {

    var info EnvOs

    cmd, err := exec.LookPath("uname")
    if err != nil {
        return info
    }

    rs, err := exec.Command(cmd, "-r").Output()
    if err == nil {
        info.Kernel = strings.TrimSpace(string(rs))
    }

    return info
}
