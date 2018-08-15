package sarsar

import (
	"time"
	"io"
	"os"
	"bufio"
	"strings"
	"fmt"
)

type sarRecord struct {
	time time.Time
	data map[string]string
}

type sarSection struct {
	records []*sarRecord
}

type sarFile struct {
	sections map[int]*sarSection
}

const (
	SECTION_CPU_UTIL                     = iota //= "CPU utilization"
	SECTION_TASK_CREATION_AND_SYS_SWITCH        //= "task creation and system switching activity"
	SECTION_SWAPPING                            //= "swapping statistics"
	SECTION_PAGING                              //= "paging statistics"
	SECTION_IO                                  //= "I/O and transfer rate statistics"
	SECTION_MEM_UTIL                            //= "memory utilization statistics"
	SECTION_SWAP_SPACE_UTIL                     //= "swap space utilization statistics"
	SECTION_HUGEPAGES_UTIL                      //= "hugepages  utilization statistics"
	SECTION_KERNEL_TABLE_STATUS                 //= "inode, file and other kernel tables status"
	SECTION_QLEN_LOADAVG                        //= "queue length and load averages"
	SECTION_TTY_DEV                             //= "TTY devices activity"
	SECTION_BLOCK_DEV                           //= "activity for each block device"
	SECTION_NETWORK_DEV                         //= "network statistics"
	SECTION_NETWORK_EDEV                        //= "statistics on failures (errors) from the network devices"
	SECTION_NETWORK_SOCK                        //= "statistics on sockets in use are reported"
	SECTION_NETWORK_SOFT                        //= "statistics about software-based network processing"
	SECTION_NETWORK_NFS                         //= "statistics about NFS client activity"
	SECTION_NETWORK_NFSD                        //= "statistics about NFS server activity"
	SECTION_END
)

var section2Name = map[int]string{
	SECTION_CPU_UTIL: "CPU util",
	SECTION_TASK_CREATION_AND_SYS_SWITCH: "Task Creation & Switch",
	SECTION_SWAPPING: "Swapping statics",
	SECTION_PAGING: "Paging statics",
	SECTION_IO: "IO statics",
	SECTION_MEM_UTIL: "Memory util",
	SECTION_SWAP_SPACE_UTIL: "Swap space util",
	SECTION_HUGEPAGES_UTIL: "Hugepages util",
	SECTION_KERNEL_TABLE_STATUS: "inode/file/kernel-tables",
	SECTION_QLEN_LOADAVG: "Queue-length & load-avg",
	SECTION_TTY_DEV: "TTY devices activity",
	SECTION_BLOCK_DEV: "Block dev activity",
	SECTION_NETWORK_DEV: "Network statistics",
	SECTION_NETWORK_EDEV: "Network device errors",
	SECTION_NETWORK_SOCK: "Network sockets",
	SECTION_NETWORK_SOFT: "Software-based network processing",
	SECTION_NETWORK_NFS: "NFS client",
	SECTION_NETWORK_NFSD: "NFS server",
}

func (s *sarFile) parseSegments(line string) (time.Time, []string, error) {
	segs := strings.Fields(line)
	if len(segs) < 3 {
		return time.Time{}, nil, fmt.Errorf("line header should have 3 segments at least, but line was \"%s\"", line)
	}

	timeStr := fmt.Sprintf("%s %s", segs[0], segs[1])
	ts, err := time.Parse("03:04:05 PM", timeStr)
	if nil != err {
		return time.Time{}, nil, fmt.Errorf("invalid timestamp: %s", timeStr)
	}

	return ts, segs[2:], nil
}

func (s *sarFile) addSection(line string) (int, []string, error) {
	_, segs, err := s.parseSegments(line)
	if nil != err {
		return 0, nil, err
	}
	if len(segs) >= 2 && "CPU" == segs[0] && "%usr" == segs[1] {
		return SECTION_CPU_UTIL, segs, nil
	}
	if len(segs) >= 2 && "proc/s" == segs[0] && "cswch/s" == segs[1] {
		return SECTION_TASK_CREATION_AND_SYS_SWITCH, segs, nil
	}
	if len(segs) >= 2 && "pswpin/s" == segs[0] && "pswpout/s" == segs[1] {
		return SECTION_SWAPPING, segs, nil
	}
	if len(segs) >= 2 && "pgpgin/s" == segs[0] && "pgpgout/s" == segs[1] {
		return SECTION_PAGING, segs, nil
	}
	if len(segs) >= 2 && "tps" == segs[0] && "rtps" == segs[1] {
		return SECTION_IO, segs, nil
	}
	if len(segs) >= 2 && "kbmemfree" == segs[0] && ("kbavail" == segs[1] || "kbmemused" == segs[1]) {
		return SECTION_MEM_UTIL, segs, nil
	}
	if len(segs) >= 2 && "kbswpfree" == segs[0] && "kbswpused" == segs[1] {
		return SECTION_SWAP_SPACE_UTIL, segs, nil
	}
	if len(segs) >= 2 && "kbhugfree" == segs[0] && "kbhugused" == segs[1] {
		return SECTION_HUGEPAGES_UTIL, segs, nil
	}
	if len(segs) >= 2 && "dentunusd" == segs[0] && "file-nr" == segs[1] {
		return SECTION_KERNEL_TABLE_STATUS, segs, nil
	}
	if len(segs) >= 2 && "runq-sz" == segs[0] && "plist-sz" == segs[1] {
		return SECTION_QLEN_LOADAVG, segs, nil
	}
	if len(segs) >= 2 && "TTY" == segs[0] && "rcvin/s" == segs[1] {
		return SECTION_TTY_DEV, segs, nil
	}
	if len(segs) >= 2 && "DEV" == segs[0] && "tps" == segs[1] {
		return SECTION_BLOCK_DEV, segs, nil
	}
	if len(segs) >= 2 && "IFACE" == segs[0] && "rxpck/s" == segs[1] {
		return SECTION_NETWORK_DEV, segs, nil
	}
	if len(segs) >= 2 && "IFACE" == segs[0] && "rxerr/s" == segs[1] {
		return SECTION_NETWORK_EDEV, segs, nil
	}
	if len(segs) >= 2 && "call/s" == segs[0] && "retrans/s" == segs[1] {
		return SECTION_NETWORK_NFS, segs, nil
	}
	if len(segs) >= 2 && "scall/s" == segs[0] && "badcall/s" == segs[1] {
		return SECTION_NETWORK_NFSD, segs, nil
	}
	if len(segs) >= 2 && "totsck" == segs[0] && "tcpsck" == segs[1] {
		return SECTION_NETWORK_SOCK, segs, nil
	}
	if len(segs) >= 2 && "CPU" == segs[0] && "total/s" == segs[1] {
		return SECTION_NETWORK_SOFT, segs, nil
	}

	return 0, nil, fmt.Errorf("unrecognized section header: \"%v\"", line)
}

func (s *sarFile) addData(sectionId int, headerSegs []string, line string) error {
	ts, segs, err := s.parseSegments(line)
	if nil != err {
		return err
	}

	if len(segs) != len(headerSegs) {
		return fmt.Errorf("data line has different segments count with header line: \"%v\"", line)
	}

	section, found := s.sections[sectionId]
	if !found {
		section = &sarSection{
			records: []*sarRecord{},
		}
		s.sections[sectionId] = section
	}

	record := &sarRecord{
		time: ts,
		data: map[string]string{},
	}
	section.records = append(section.records, record)

	for idx := range segs {
		record.data[headerSegs[idx]] = segs[idx]
	}

	return nil
}

func parseSarFile(path string) (*sarFile, error) {
	f, err := os.Open(path)
	if nil != err {
		return nil, err
	}
	defer f.Close()

	sarFile := &sarFile{
		sections: map[int]*sarSection{},
	}

	buf := bufio.NewReader(f)
	sectionBegins := false
	lastSection := 0
	var lastSectionHeaderSegs []string
	for {
		bs, _, err := buf.ReadLine()
		if nil != err {
			if io.EOF == err {
				return sarFile, nil
			}
			return nil, err
		}
		line := string(bs)

		//file begin
		if strings.HasPrefix(line, "Linux ") {
			continue
		}

		//ignore averages
		if strings.HasPrefix(line, "Average: ") {
			continue
		}

		if "" == line {
			sectionBegins = true
			continue
		}

		if sectionBegins {
			sectionBegins = false
			lastSection, lastSectionHeaderSegs, err = sarFile.addSection(line)
			if nil != err {
				return nil, err
			}
		} else {
			if err := sarFile.addData(lastSection, lastSectionHeaderSegs, line); nil != err {
				return nil, err
			}
		}
	}
}
