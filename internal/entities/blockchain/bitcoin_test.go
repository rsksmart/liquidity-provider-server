package blockchain_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

var p2pkhTestnetAddresses = []string{
	"mzVTV2cwEpBLWcsPckthMA7NZ2Cw8ojWf8",
	"moAJa2o3ggwSweWP8qmPznML9bRke9vsNc",
	"mtEjnJR9MPQEkpqvuKemsFHdbL1cXrmEnK",
	"mucJ4ZUriexq6poQv717ArtkJxEEZXVLMe",
	"mtcy9BuREZ54EbArDTdTz4we6M7AYybU3P",
	"mtwn2DFCfiJ75ppTYvAULs1v3nqyKo5Dnn",
	"mpihFFPkg7HL1RbQcfogE3pksUqMRVjsqc",
	"mseeU5gbFbEuDpoZQ5fwBi5UyTUmE3wD3j",
	"mncLnq3YFdf8vvRkLySmJhX4YuVttuAp9C",
	"mrQ6EZRMgKdSXQN4rK7rNxwqXV3Cnjupb4",
	"munRXv4sfA4E6PazLzgEFF2mWo19YJQbTD",
	"myTJXHz7JUcV1vq9Y9vhoqmdJaSYWfgkUh",
	"mjqo4zmzzybYMcm5HHdhrWW7N6XfHR1qDz",
	"mgDUopmALewBcNRoTCDNgTF75cFkQuCDdP",
	"mv883zByNTJmaAVZ9RbcUEMFFmGwDBvnC5",
	"misG6J9JgskhKsf9PfrAagan8fx2UUCErF",
	"mj2ucWCL3SSkbJRoLExMQZPMH7fEvoejH9",
	"mxa17v5jv6jT8YkX8GVPNkt4DGzRG15ebi",
	"mq4LomGAruBzQm3UcM9Seu6TFtCSsXX6zx",
	"n3KxPinGaXMLxJoVBiaaQa4N8x1i2aHNH1",
	"mvoBtYKhVrnvWWnqVEGn7awGNRisw4PkVF",
	"mxZ7CYDMUaKZt8VVsMvoXif3vRm4aK9FGD",
	"mmWHzFtKBfBnitWBHX9VEJee9skear1qZX",
	"mqfkut8yAQepC6kj8HtiKXPC66MefCj2Yc",
	"miE2R7EWfhJ7mMfSems4moz8HejSR3igP1",
	"mxTYtVFuZxtE6x19VAQ1bZEEBrFAe1NJfS",
	"mxkZ6CK1edZZuEEwngy52Tm5WxLRaHXKrY",
	"n4XLcxmXTXdhEQaBDGPP9FU8CReaVcNVQm",
	"mjVA4tXotqXyP5qw3d6iNFSJfVmkPLThaW",
	"mxgUSkFfUZPmmQRmhNtb7s3uTWFdRzuvWQ",
	"n2Y3y2ZkjoDPNGM3GYAcgX91na35Prrzr3",
	"mfXSEKxMzPkejLwrZeNzxE72QnoC5fmELi",
	"mx17QNaAh9RCSQ9zajyaSN4dfjNjPszjUf",
	"mzgWf4TVPbjYvDAKRQdWHEDyZFQCEHHrWh",
	"moig2R8WViDPbEsCXYSCBQULgEVHs93eeS",
	"mgFWCp1UKkWt5SFQ5MAa37Fd1JoYTGwmCK",
	"mjjo6muRiQzuzqvpTAK3A4odD2RXjHLTz9",
	"mos8ZoeAZVpJSafT7qgv3Pf7Q3WtKXP1E7",
	"mtd3ueANi3hsU1NVVdoUnuhSwy5LdzVXXr",
	"mw1GKH6wpUmu1w83iP93asovyNod11ByrQ",
	"myVVnWrgcuqmGj7f6CD3EyCLCi73sjjFLD",
	"mg4pSPq9azViCbfSKpzZFHsuG17xG4pn3S",
	"mi9fMknNV373AsJuFPKrQgqaMopnQsqZft",
	"mrcAitXFKFcGWZUPTRnDv7qodxq3rAH2C7",
	"mtwdVmgAXE5vfk1QziXnMhJb5DG5icognD",
	"mubMX7zb8gmNiSd8qqB3NQSoQ8dc17qkkt",
	"mrE5MMQbEbcEQQeGiMKzadd6h79XKtycVc",
	"mzdBTxxvDkKHye3ZbUqgxBnMTzTxN3Mv9r",
	"mzvKAkcbBpGdDXdA3xCJEgtGLmb24mJ6CS",
	"muoNEVG8kscidB8ijrgsY9HGQMfP1FWUbH",
	"my9bJuzTzaudoUA3jYZAisjDMpmigwFrL7",
}

var p2pkhMainnetAddresses = []string{
	"1Ld4rjLYfmtW2HrF47mHjVLfaU82ku6GFX",
	"14KtqniYfjBw6BxyGadNqJRE2TfrZ2iDex",
	"1LdzDJ51dpt2cg8jmcaay4ypVYDykYuRZc",
	"16Pg2ctWtwbAGvZVNnQCEpWyLV85BwEDcX",
	"1Ab1Jfe6xQHzL8RHoHDukDQBEks35KFWHC",
	"19hiJTQpZyT3C7Hu29dJE2YYToCeKp6cGu",
	"1LQdpgVCY2nYzsoRNRHWhuCLxMpYzb6zzg",
	"17oKLsbZsd2BZdCDn1dbrbk2TT9HSzw2aM",
	"1NbJonAytRKfCFkvGcQNEUCXAFnf17bYQG",
	"1MwygkmvJHwwG934EbtkjhRUFyfMHLEPi9",
	"124qYovdcDQQzBETVsyKBA8vYPpQTjYo7n",
	"1BMjg4kirMSZ49h5DbR4fo6ZrAT6CaeQSc",
	"14uFoAvDGsa25VSYk4AtehxcGfhTLscyGB",
	"1MKLYiRiJ94p7T5BVi7LjkxQsXZXmfaDNY",
	"19QVu5DkLHhvikdkCF7AHCxQK89xQFJGj",
	"1PWFvoD73rsbsTqojvXGM5m7xHdM9iYDSM",
	"1PC1Vz6zEocJdWvYe9rwBpDH6cJ23Jf57g",
	"15PiaAwCXC6R3Ub9w4CjmWqC537C1GA57Z",
	"17Q7ZhQxryfvhEhBv3tXtXrWhczCfho8Vg",
	"17KXpqaBazhkNXxM3swFL1zkTzmKJFCmhX",
	"1DLaPDD6tRWzqbvt9SQ6cdYcTJxdYmqpmG",
	"1GzfGiV8T8MTADzrv1XBtwt6GDFNUNUbzR",
	"1BGbjogTCqUwB7MmRBjBZob1wydau8BE9f",
	"164XxTfzPyLY5DDXYfS8sCnChoeHZBNT3m",
	"1Dn66Pqx8p7Ky8gc9G4x6x6riFZzZ4jJqF",
	"18FixbDw6EEPqEGXphdw3BaXHb7xUYx3ff",
	"149RXsdcY6s3Vw8rzfjg2pTNkCGWqhF1dn",
	"1QKLNCmyqVPUZneefDw3Yj133Yf47E53My",
	"15u1itBJhUJTY9ZmqpHJagBDuoN2fBrfbK",
	"1MK3xXgsSzwSY5V78fiRX5PCPWD6iVg9ry",
	"1A713zQYJLoQKKZLxyqawfdKwpSuy6Cu19",
	"1Gkd5uTfvYjpubq1ZaYCeY8xDuVAcUoVnF",
	"14JuvbQg2XJztYpyauX5WTsHdQwp1iX5pa",
	"12YkAVgWxoJF1zMYJuCjbrESNWZY8iyNPC",
	"1D6tCRyqecKUa2xGFH5Qqnnd3WLrC8JN8V",
	"1LS6ZSCEzhfmP1g7MEhjson7TKCq1RJ5eP",
	"147gsg9Y4npuDrubtBVFYxGfcxCcnhfXAu",
	"19fkgqMpCkRJF6kXwyeT3FLLpBsN26Thad",
	"1NqeTZNq8P7P8755UnPC782zKPhrgYtyin",
}

var p2shTestnetAddresses = []string{
	"2Mssid57MEEAzrjiVCG87PsRDyW52Wnjwso",
	"2N7A51EjcWVSbP3ag7b48oJ8QiEm5wk3CmC",
	"2NDcdkBpUneEgKvcmYpC4n9faoRiALeEzbP",
	"2N5WuV8JPs4sjBkjQsFkjTDPGDgjBhWxXMj",
	"2MwgCpDmUiHMUsqmxqQmEM3AEoeusXxver8",
	"2Mu1WgnUKR2C1R7FuVbfC2jRoFbAtcJCtP7",
	"2MyKtiQcyAgQQdPDmBJcT4UMn6jyVdKVwxg",
	"2N9TujzoGXb4LkzCQscMuuRSwaTVPhFfYSy",
	"2NAvNk9t8fKwFRhMJGgY2wMLRboN8DHXeBB",
	"2MzqFJDg3QcbzeWX87XpxRpHZm6SSNoGdoF",
	"2MuBrh282woZkhycbpKAZ8zEptTAEtRSM62",
	"2NAobddanbJX2AQmgv8iSCHTTT3LEFPkdN8",
	"2NEwLSjVVuDfr6mpvkbVNy7MpoxHhLWduvy",
	"2NAgRDmhg9zDeJQLFx5tfGhzs2zZfLR4L9o",
	"2N3s7K9314WDqSNS8NKggcBpgAhsVBjoR72",
	"2N1h7tqTub3GevTpdGnkyBYChApwN25w9Ex",
	"2MwnxYx5vvCj46ThT18aB1e5K3rAhyVh8wz",
	"2MzMzYAHzbTHFg4VEoidRjkuVDuFWce4vm4",
	"2MzDws1oDEyG2f6hkubUXCUUnu7aLicXtEN",
	"2N34YniLPdNJ33vBSbzUphmhM5F6XcX2vgS",
	"2Mx4KfFN8vspsXkEZZGpZHdmjkTjnWQ4SCU",
	"2MwGLNdaZSkVvcqpuJGo17AonGGWHfrjF1N",
	"2N3vaXgbkvb4Vv3etHq45KQzCM2CczyTJGe",
	"2MxXigcEyqQaErgA9qCjuyHg3c2YYqEY8dd",
	"2MtS6SDDB1uN7zb8fTFhaiCfJQWLBQYApg1",
	"2N6PiKxASAKihec1Q2xLHo8stXccCZsFJg8",
	"2N4oTNvcvttro4chtve8qLCmEnp8cbt2sPR",
	"2NBPF78FpdLBDxu6FvmZnw7K5462yn3SeVN",
	"2NE73iNV74T4FTiULnRkwaAUxXKQugBeXDK",
	"2MwgNDvWAnjsMZ2Jp73VhdyaFc38Kx4rFMP",
	"2N8EWGy86QJt67m57DbRBN9V5Ksj8m6ZL1E",
	"2N7QybDvzGcyKDJXYYkfzvFWwvHyrpc1M9S",
	"2NBdfBtz7Wfht4QTAuEWCdyTnm3HucPbPU6",
	"2MxeuFGBzQrA54YYEPcVMFaXKYVfBTBBKB3",
	"2N58BH8rEq9Ku7HuJbZvKX6WRywdNmoVrnA",
	"2NEv9myWjPZVLxEUKEMC5vbKP5euMTxXZyM",
	"2NABLwMDsitb13Y4Peuoy5M5fishWaFKaV8",
	"2N4qmbZNDMyHDBEBKTCP218HV1LhxCMRMti",
	"2NDt1sRinReq2qdJf1wTojfhx7WgdvyumEB",
	"2N8XxpAd9teiPEeBgLfi6vuAgq4KA4CeUaz",
	"2NDTyc7YZ2Z3nnFqswcF1YjnuLAMjxnpiBB",
	"2NEEHG27iCKuBWSmCKXygXv4uZ1yPWu6vfL",
	"2NBnogvRFVas4q8juTCohgv7qFiKPomTvyr",
	"2N2ZEQTGFeTDgzHoPX1KBAxejpWCwvhF61f",
	"2MupuVBvFvcyEnGdrRBixhZrTxRxctaQHUM",
	"2MwysnGXHByGyvyJKkwwdobdoyhcCC2LKkF",
	"2N8Kr3eyzSwEPNi8DR1GEhuwnvsQJCvquXP",
	"2NBTmXUJkgx5SFsiDUffNaFGBVaC2W6GjxY",
	"2NEe5WztmDRC8TvUo2bqurK1dNMRvPf52qQ",
	"2MzBVRSRJqmGhsYXc35ecWKyLw5p3Y7bdLs",
	"2Mv4aWbKBc3k9ZxF8SN2k3VMuAWfU3Fv8hB",
}

var p2shMainnetAddresses = []string{
	"3HwGVBDv7jPuS8ncTFm36UW69h36BsdjqL",
	"384B1WkNYMm58o6F7w1yVEZwX48EvFKbSa",
	"3PUWoEeffv7y1WeYEF9v2gA9AxQ9v1FVLn",
	"3KXu18nZ5qE9mwap3sKi35Du9rdSn8hHuP",
	"3MnspZrXnAvYhCrGLYNvRreHiHsz587Z5q",
	"35Vebz7gEQk8pEpi5FyZy4fTGngBificfL",
	"39NNxXAg35E256YeU3ofnPXzkq3aAoMJdx",
	"3Pu2Jp63S5WfFTE3bZM9kPPUFAqQvXCEBZ",
	"3QD7radHE3LUjWQGtpXvCiZJeGmPPLfy55",
	"36xT5Ly3hyPATenoFvXn9vD21gLFecjSpc",
	"3End97i9TYEkSFBEpAP6mKMYHDJHcKAJ2S",
	"3JRL2REh3KGJaKLNuPXCz4K2ZSCVbzhBKh",
	"37ofEuKT4cZhc1iZY5k1uRx1ShoRwtXmXi",
	"3DEAo1q2dABPzoxvcqsEXYSzDBvBHpyPkp",
	"3QmBRzpsUDxUjt7TzNGRcHHgZRn3gxszyb",
	"31tvn8gdKgg525KNZM4wevi87GafSzD1FU",
	"33BTakydSPnJfSfR13foniEsPCB2nuHiCb",
	"3Fi3ywSD7eBEtnoUuiW7zFSBmpATd3YDLs",
	"3DCT2YtzwZZYdr3pPEhqVjA1Amak89YrHf",
	"3EuQdJ651cN7Cv9jJk2EJPdRhKT9JJFpt8",
	"3BcLTSd24JRtJhLcKqkeF83rFmFxxY5qH9",
	"35ficBSmjzqdSqAWCKgZCeeTG5q5mB5EFj",
	"3A2YVYwmodEATGGKSjJ9vx7jY5YLkq9kst",
	"34JdvFMpFfBSk6ibR1HAzjTVaZ8Dbt1YuL",
	"33tB53JpXTNeDvMUeXKLU3VfYMLMc6oLab",
}

var p2wpkhTestnetAddresses = []string{
	"tb1q2hxr4x5g4grwwrerf3y4tge776hmuw0wnh5vrd",
	"tb1qj9g0zjrj5r872hkkvxcedr3l504z50ayy5ercl",
	"tb1qpkv0lra0nz68ge5lzjjt6urdz2ejx8x4e9ell3",
	"tb1qqgzlw8yhyj6tmutat0u5n3dnxm3y6xnjp53wy9",
	"tb1qcc4j0tdu3lwfl05her3crlnvtqvltt90n5s5m0",
	"tb1q7k3nex0gssyucqvz7xpk25wzqfpc56ve9myzqs",
	"tb1qur2ztvmx4tqdxa35js04zuqhwx624z3nyuv97l",
	"tb1q3n9lhc63xwkfrj25sy2gqf06r77vzqcryxqe7k",
	"tb1qc706vv5vyqqz3drx080c3um3t2ylze8fuxuujd",
	"tb1qj36sglrgm590mdkht0dququug73azfxcxdnhk7",
	"tb1qu05qzyrlqaeth5j0p0fkxek4tp5huffveryewa",
	"tb1qv87afr2gu8g57v7u39g6h8txwx7afvkzc5ja0j",
	"tb1qc85tj3uc7auyndw4en02vk74ty2720rpppaxd5",
	"tb1q94u4lqykk9m387p59sqcvks7dhjpaz8tf5kdkp",
	"tb1qclumyllxep6gp9rnv7wzks869t7w5ct9rznuzd",
	"tb1qv55ksu4ll2xmekru50nknac6zkq9c87mf387n2",
	"tb1qw853dsyg9745dm5q39zmgnk3m6ldr879q5rht6",
	"tb1qtzx5vjl37rl8nefn4ppdqvwrxw9cvrqf2w632d",
	"tb1qxylu8c2r9ypucc7jas47n3dy400kx0n4hd5g2t",
	"tb1qmznvxadpzmzc5x3q9cvwsu0vrud967hnryndq0",
	"tb1q700p8wdp9t6z3f59f009uwqkf8nct49arn9zh4",
	"tb1q0dgpnxve9utzc3m38zmw7drh0ekdw20gup6sqp",
	"tb1q36tpm7eu706v0ut0hap6yjuehgsg53rg280tc9",
	"tb1qwnmmmrrr7hw60yulw2rx50ne2tkktj729076zf",
	"tb1qag9uv7n266eyf6d88xc3e5nmek8sqe6aqxmfpp",
	"tb1q22cm3qarlpj3gnf5h03kpdhaftdvf98q58dp75",
	"tb1qja7532egus56jkjnu6xgf9nh96q9up7gq5473m",
	"tb1qug3kle73ze6wcstdc4wunkjxapqnaeetprqjql",
	"tb1q66e97gspk233et7k24334zm2femvf5tpsq8ggm",
	"tb1qjtl57d37ccadme30hv3jhytt9gc9p4dq9zrz49",
}

var p2wpkhMainnetAddresses = []string{
	"bc1qg5d579rlqmfekwx3m85a2sr8gy2s5dwfjj2lun",
	"bc1qtqxd29s9k3tj3rq9fzj7mnjknvlqzy8hsuzs5x",
	"bc1qv245zr29zw5urv5fy00c6km09l302fmlftf0aj",
	"bc1qw4z64jjvuxyddjdcm88yt0ln7fntkyw0w6wqhp",
	"bc1q8d7e3jrhsf8tj9q28x3msf8c644hdaetpqy7t4",
	"bc1qzs5h0w6zjk3gej89gz3fjqx3xv9kvngntam76g",
	"bc1qtj7y3xapmn38gra6jd3x6ua905j72plz5rv0kq",
	"bc1qx2xgfynw4fjtm8c2glv9ceue5j0k42fdp6u8vz",
	"bc1qtjrcft537z0kwcl7efucqh3f6xyvkqc908g9ve",
	"bc1q4a2le8yt6l6x0t9cfxmmx35runhf58lktvq8wn",
	"bc1q5qxpf5ca4s9k420h6vx9gegdqtq77sxdmqxca6",
	"bc1q9wjmgcfny0rsyhd727y6t3xv8wf69ggmk04msw",
	"bc1qufr3cd7kmr9kufd62wh3d0jq44zqrm4yqjyc04",
	"bc1q455afukjmkm0v9fpldqacc8jdevfvcylxtz97y",
	"bc1qyg6m5eyjlder7cg8ja5lw788jn9pfsx4ypyr9t",
	"bc1q9qfptdws3l0qqrhc8fvfezw3uv3f0vtrzz4m57",
	"bc1q7uq8kxtnu6ya8l3k7c4w7avxwweqa8c2nfv7zh",
	"bc1quk39lyk9tya3z34g7zarfnxegv855584pmy9ud",
	"bc1qfy323zdtp7dd6pjxkz5k46l7v6kktmsp2vrsfd",
	"bc1qkmm0yj2heftt0safn7cgqrexdzhujz53kfpq6w",
	"bc1qwr5glcl76g7hwnx8kunqsxuapm6vq3mg8yel50",
	"bc1q8eg9j6khqeqhrjsmc5890nnpjyaj3wdkkq6485",
	"bc1q2rrlg43vv5snstv3mvc79mfr9amfw4yknhjtew",
	"bc1q9da5fx8eerg4m40vkqc2mm24nykxhtpw9sfcw3",
	"bc1qma0pelvcshhq59wfur9p5rhacjyk0lmfdr53vs",
	"bc1qzulaxy8fmvk8a92sec8s8u0xcqwcxw4fx037d8",
	"bc1q8sr9tv9ng4yd8s6s9eenfs7mh24jv64vnwzl0p",
	"bc1q5pfzfxmtx3kn7j8wqwe6336tmg0n5lmpqss9kx",
	"bc1qgq506g46u2dnua70k3dypu6r7xu3kfqeee3c38",
	"bc1qk2r5qt94fluyehjhr6neka0agpxung28pndjly",
}

var taprootMainnetAddresses = []string{
	"bc1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg3297",
	"bc1py8g4v4ehll399qlpaxyxykg37pszhad9yg0dphxvjhdmhy7f08vsn43s6p",
	"bc1pc09cafvlgu5ykmxyyzr4gu5qwx9a2zz6fz3lljeyddc9z7n75n9qfz7ckr",
	"bc1p74k39706fe6n0qv5k30z4xpqd8gcf8apyzn9s5rujkz4jln2u3fqwwta94",
	"bc1petgnkphl82md05d84gwee0alkuzpphfjy8ycxs932ngvdx8z8u0s3dwj5t",
}

var taprootTestnetAddresses = []string{
	"tb1pqaas5xm75dny58s452949c9ak5qd53shfkln490ju4ny2afs2ldsput844",
	"tb1p25h0xs3840q7aex3kl9dshd8q99qzaxkh8r5p70z54r4ykmn2rtsgcsj34",
	"tb1p7hvw8mnqlrtp7ffa8wzmhq7vddegffdeus4sl0yj6fw54zjda36qhc5q8y",
	"tb1pnqdr56lugmtrcxtae8k9cfe7hve8986ud0daktljsh93wf8q7u4qhc2q3c",
	"tb1pa54gmj3dzr9g5p7qx6kupqg9xkvtv2cdcty78wgyaycxtqc72h5qlqgz2c",
	"tb1p8lkxfnps5wd6rsrusvytp8zllrmxz05e0ttessnhyzwl0kusc2as4s72wz",
}

var p2wshMainnetAddresses = []string{
	"bc1qhnumvtg3c9xj2q7jmt8xnk4p5kmk52ffqwax8crfn4hqtry6qseq8vahua",
	"bc1qv47nn097m6hujqadw6kgt5hsk9h06k7tgq05empl3nn3mska8cfqpkjl36",
	"bc1qj8pqhwkv0k6h2tm3wtqu793njkvfd66dva04zldpdcey4sak5h3qx3n8nz",
	"bc1qem2ta6uk98rfr779t4wftq4qjtr3xtja9vf9yy3rgtczapc78j3sxa6570",
	"bc1qazm8jprsdjxn0qq77yrzw7m2340ys0kuuylg05vul4t5ll2lhduquuhngw",
}
var p2wshTestnetAddresses = []string{
	"tb1qzda4qlkdpjgmwxt9zr29pphhzqf2ku09p7dj33qyugqn80kg5muq8x0wyv",
	"tb1qgpgtqj68zwsdz7xmvqxxxaan7dcfgu76jz0cfzynqgrtvdsxlyqsf7dfz8",
	"tb1qkp4lxc09e34cc5vw383j42rgacurp7wrpnwjmvazv6g23c2ydz3qx5tfhl",
	"tb1qzhu8fjgw5aaqgv0q2jey4dnc3pgcr4cks858d6eaf97ljxywe70qwwsdku",
	"tb1qzda4qlkdpjgmwxt9zr29pphhzqf2ku09p7dj33qyugqn80kg5muq8x0wyv",
}

var p2wshRegtestAddresses = []string{
	"bcrt1qtmm4qallkmnd2vl5y3w3an3uvq6w5v2ahqvfqm0mfxny8cnsdrashv8fsr",
	"bcrt1q93c4pfcz8lye3d0j9y7dp9sr48rj5kv34gfds4f0q9nmwd6afhrsy5hv6e",
	"bcrt1q7fs9g5c3kp3r9c6e8qx0y6lfjc5z7xq0xvd3crw0pqkz9m7v84qq2v9l6q",
	"bcrt1qdpszwqsc2vnf9sh4jl8gq0g2qpj3eyg89szukln6wxdtr3w6suvq3hyt4q",
	"bcrt1q94s5pfq9eycr98rdscpjxzwc7vzh8nfrlx2wzt8hztpyzvcajndszww3y9",
	"bcrt1qzv8ykc6vcq4fzzklgk4e3wts3r6fzjd7jqh4k3lw6tztkgv5rxmss0q6n4",
	"bcrt1qv95xqjc42ngc60gptf5pv5g4gq3hyuqytxvhdjq4nwvxd93w5lzst3t6ez",
	"bcrt1q96s6z5fnr5vxhrgslh4gk5g5vj95w6krj3jsl8r60lm3fvkljqdsg3ch8v",
	"bcrt1qxyprwrln9w5qk2ghnys5xfwz0s3u8jgrdm2vj42y27tqeqp8zznq0kwkcn",
	"bcrt1q7yjkl6r5mg8j9hyt8gx6f2h5jq4kt8ks7w0q3r4wlcck9yks4vzsryvw8r",
}

var p2wpkhRegtestAddresses = []string{
	"bcrt1qg59jse8wha2z07w8jsjd3my7530tk0vqh7hv09",
	"bcrt1qpncuuk7mflzwyfrv7uwtskl35qdxs0uuwzwhvf",
	"bcrt1q32chehq23qgu24h5jvy4407ypvmrfvut7vssyg",
	"bcrt1q5rg6mxr9z4myvs6c8gkh7lcrcj9eexava5979n",
	"bcrt1qpq0hdfj3fwtusy3vjmhkmz9y7epvveqalgt4te",
	"bcrt1qp6u0t6evapxq7q89yv7xfsq978dwnxeskkup7t",
	"bcrt1qkeya2ky446gl6kqm86nj96tjxhhsf9ctsha9xj",
	"bcrt1q9pt3z976daw3jw3jncqs83vys6ay567hdwhqgt",
	"bcrt1qvgfxpgzhhd6ktq0rfc5lm8ysu653uyxmqglvtf",
	"bcrt1q4lqflp2gywyuhdudawm4kncq2m2llaydh3srju",
}

var p2trRegtestAddresses = []string{
	"bcrt1p5chgc8qnl250kjp4956u690hqmkr7px2xdxq0hnrd3e69slumsfsj7wslt",
	"bcrt1pmxnjedrmj4l7muyezf62x6cw6e385a36v9sk9znlxz54dwwtldlqxm6p8h",
	"bcrt1pgtyx9g4xw9kx3m6v9fkpadnrp820t5mk6v4vry4sqadkg6fehp0qqtq4hp",
	"bcrt1plmdsnxc9mmjt4dks4qapzd8cr2m407excd0c4a3zjt6lchzlhd5qtx4r95",
	"bcrt1pt7qtpvcevgkl4453nk3ctlyc89qpurjyvakpqwu5l3prsf7jlynsfzj337",
	"bcrt1p62ulphlvyrkfxjy5xun02qshdmdsmdxrssa9fltgc9pranxdpceq0zxehk",
	"bcrt1p7tpuqkurf0z4rq0auepqqnlycp454mtcvmxks5wqpp5dxpgxw64sr9gfl2",
	"bcrt1pujvqun9necxr76k235kme0auxx2wc5prhapa303rk6kwt7l4wvtq54cpsq",
	"bcrt1pr9mhhhd706m6gsyddfpkp3jhpj4cneqrt8sy7whm8rlts8t2cs8s2dk38m",
	"bcrt1p7m6rue7kxpaqc56rcvym5n9jwtumsvjv00qqk0u4c953tps0t0sqklkgge",
}

func TestBitcoinTransactionInformation_AmountToAddress(t *testing.T) {
	address := "2N2Sg8C2uX1YtugYSxEQvRqf9V2EivxcWER"
	cases := test.Table[blockchain.BitcoinTransactionInformation, *entities.Wei]{
		{Value: blockchain.BitcoinTransactionInformation{
			Hash:          "0x1234",
			Confirmations: 1,
			Outputs:       map[string][]*entities.Wei{address: {entities.NewWei(500)}},
		},
			Result: entities.NewWei(500),
		},
		{Value: blockchain.BitcoinTransactionInformation{
			Hash:          "0x1234",
			Confirmations: 1,
			Outputs: map[string][]*entities.Wei{
				"2N1nBfGejU5iLEqAS42fBKJ1Dw6mw4su8eQ": {entities.NewWei(100)},
				address:                               {entities.NewWei(500)},
				"2MvHto2NWaAtiMeDsy2oAHesnK8Rug3Lavc": {entities.NewWei(300)},
			},
		},
			Result: entities.NewWei(500),
		},
		{Value: blockchain.BitcoinTransactionInformation{
			Hash:          "0x1234",
			Confirmations: 1,
			Outputs: map[string][]*entities.Wei{
				"2N1nBfGejU5iLEqAS42fBKJ1Dw6mw4su8eQ": {entities.NewWei(100)},
				address:                               {entities.NewWei(500), entities.NewWei(1100)},
				"2MvHto2NWaAtiMeDsy2oAHesnK8Rug3Lavc": {entities.NewWei(300)},
			},
		},
			Result: entities.NewWei(1600),
		},
		{Value: blockchain.BitcoinTransactionInformation{
			Hash:          "0x1234",
			Confirmations: 1,
			Outputs:       map[string][]*entities.Wei{address: {entities.NewWei(400), entities.NewWei(1100)}},
		},
			Result: entities.NewWei(1500),
		},
		{Value: blockchain.BitcoinTransactionInformation{
			Hash:          "0x1234",
			Confirmations: 1,
			Outputs:       map[string][]*entities.Wei{},
		},
			Result: entities.NewWei(0),
		},
		{Value: blockchain.BitcoinTransactionInformation{
			Hash:          "0x1234",
			Confirmations: 1,
			Outputs:       map[string][]*entities.Wei{"2MvHto2NWaAtiMeDsy2oAHesnK8Rug3Lavc": {entities.NewWei(400), entities.NewWei(1100)}},
		},
			Result: entities.NewWei(0),
		},
	}

	test.RunTable(t, cases, func(value blockchain.BitcoinTransactionInformation) *entities.Wei {
		return value.AmountToAddress(address)
	})
}

func TestIsSupportedBtcAddress(t *testing.T) {
	var supported []string
	supported = append(supported, p2pkhTestnetAddresses...)
	supported = append(supported, p2pkhMainnetAddresses...)
	supported = append(supported, p2shTestnetAddresses...)
	supported = append(supported, p2shMainnetAddresses...)
	supported = append(supported, p2wpkhTestnetAddresses...)
	supported = append(supported, p2wpkhMainnetAddresses...)
	supported = append(supported, p2wshTestnetAddresses...)
	supported = append(supported, p2wshMainnetAddresses...)
	supported = append(supported, taprootTestnetAddresses...)
	supported = append(supported, taprootMainnetAddresses...)
	supported = append(supported, p2wshRegtestAddresses...)
	supported = append(supported, p2wpkhRegtestAddresses...)
	supported = append(supported, p2trRegtestAddresses...)

	for _, address := range supported {
		assert.Truef(t, blockchain.IsSupportedBtcAddress(address), "Address %s should be supported", address)
	}
}

func TestIsMainnetBtcAddress(t *testing.T) {
	var mainnetAddresses []string
	mainnetAddresses = append(mainnetAddresses, p2pkhMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, p2shMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, p2wpkhMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, p2wshMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, taprootMainnetAddresses...)

	var testnetAddresses []string
	testnetAddresses = append(testnetAddresses, p2pkhTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, p2shTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, p2wpkhTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, p2wshTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, taprootTestnetAddresses...)

	var regtestAddresses []string
	regtestAddresses = append(regtestAddresses, p2wshRegtestAddresses...)
	regtestAddresses = append(regtestAddresses, p2wpkhRegtestAddresses...)
	regtestAddresses = append(regtestAddresses, p2trRegtestAddresses...)

	for _, address := range mainnetAddresses {
		assert.Truef(t, blockchain.IsMainnetBtcAddress(address), "IsMainnetBtcAddress should return true for address %s", address)
	}
	for _, address := range testnetAddresses {
		assert.Falsef(t, blockchain.IsMainnetBtcAddress(address), "IsMainnetBtcAddress should return false for address %s", address)
	}
	for _, address := range regtestAddresses {
		assert.Falsef(t, blockchain.IsMainnetBtcAddress(address), "IsMainnetBtcAddress should return false for address %s", address)
	}
}

func TestIsTestnetBtcAddress(t *testing.T) {
	var mainnetAddresses []string
	mainnetAddresses = append(mainnetAddresses, p2pkhMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, p2shMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, p2wpkhMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, p2wshMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, taprootMainnetAddresses...)

	var testnetAddresses []string
	testnetAddresses = append(testnetAddresses, p2pkhTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, p2shTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, p2wpkhTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, p2wshTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, taprootTestnetAddresses...)

	var regtestAddresses []string
	regtestAddresses = append(regtestAddresses, p2wshRegtestAddresses...)
	regtestAddresses = append(regtestAddresses, p2wpkhRegtestAddresses...)
	regtestAddresses = append(regtestAddresses, p2trRegtestAddresses...)

	for _, address := range mainnetAddresses {
		assert.Falsef(t, blockchain.IsTestnetBtcAddress(address), "IsTestnetBtcAddress should return false for address %s", address)
	}
	for _, address := range testnetAddresses {
		assert.Truef(t, blockchain.IsTestnetBtcAddress(address), "IsTestnetBtcAddress should return true for address %s", address)
	}
	for _, address := range regtestAddresses {
		assert.Falsef(t, blockchain.IsTestnetBtcAddress(address), "IsMainnetBtcAddress should return false for address %s", address)
	}
}

func TestIsRegtestBtcAddress(t *testing.T) {
	var mainnetAddresses []string
	mainnetAddresses = append(mainnetAddresses, p2pkhMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, p2shMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, p2wpkhMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, p2wshMainnetAddresses...)
	mainnetAddresses = append(mainnetAddresses, taprootMainnetAddresses...)

	var testnetAddresses []string
	testnetAddresses = append(testnetAddresses, p2wpkhTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, p2wshTestnetAddresses...)
	testnetAddresses = append(testnetAddresses, taprootTestnetAddresses...)

	var regtestAddresses []string
	regtestAddresses = append(regtestAddresses, p2wshRegtestAddresses...)
	regtestAddresses = append(regtestAddresses, p2wpkhRegtestAddresses...)
	regtestAddresses = append(regtestAddresses, p2trRegtestAddresses...)
	regtestAddresses = append(regtestAddresses, p2pkhTestnetAddresses...)
	regtestAddresses = append(regtestAddresses, p2shTestnetAddresses...)

	for _, address := range mainnetAddresses {
		assert.Falsef(t, blockchain.IsRegtestBtcAddress(address), "IsRegtestBtcAddress should return false for address %s", address)
	}
	for _, address := range testnetAddresses {
		assert.Falsef(t, blockchain.IsRegtestBtcAddress(address), "IsRegtestBtcAddress should return false for address %s", address)
	}
	for _, address := range regtestAddresses {
		assert.Truef(t, blockchain.IsRegtestBtcAddress(address), "IsRegtestBtcAddress should return true for address %s", address)
	}
}

func TestIsP2PKHAddress(t *testing.T) {
	var p2pkhAddresses []string
	p2pkhAddresses = append(p2pkhAddresses, p2pkhMainnetAddresses...)
	p2pkhAddresses = append(p2pkhAddresses, p2pkhTestnetAddresses...)
	var others []string
	others = append(others, p2shMainnetAddresses...)
	others = append(others, p2shTestnetAddresses...)
	others = append(others, p2wpkhMainnetAddresses...)
	others = append(others, p2wpkhTestnetAddresses...)
	others = append(others, p2wshMainnetAddresses...)
	others = append(others, p2wshTestnetAddresses...)
	others = append(others, taprootMainnetAddresses...)
	others = append(others, taprootTestnetAddresses...)
	others = append(others, p2wshRegtestAddresses...)
	others = append(others, p2wpkhRegtestAddresses...)
	others = append(others, p2trRegtestAddresses...)

	for _, address := range p2pkhAddresses {
		assert.Truef(t, blockchain.IsBtcP2PKHAddress(address), "IsP2PKHAddress should return true for address %s", address)
	}
	for _, address := range others {
		assert.Falsef(t, blockchain.IsBtcP2PKHAddress(address), "IsP2PKHAddress should return false for address %s", address)
	}
}

func TestIsP2SHAddress(t *testing.T) {
	var p2shAddresses []string
	p2shAddresses = append(p2shAddresses, p2shMainnetAddresses...)
	p2shAddresses = append(p2shAddresses, p2shTestnetAddresses...)
	var others []string
	others = append(others, p2pkhMainnetAddresses...)
	others = append(others, p2pkhTestnetAddresses...)
	others = append(others, p2wpkhMainnetAddresses...)
	others = append(others, p2wpkhTestnetAddresses...)
	others = append(others, p2wshMainnetAddresses...)
	others = append(others, p2wshTestnetAddresses...)
	others = append(others, taprootMainnetAddresses...)
	others = append(others, taprootTestnetAddresses...)
	others = append(others, p2wshRegtestAddresses...)
	others = append(others, p2wpkhRegtestAddresses...)
	others = append(others, p2trRegtestAddresses...)

	for _, address := range p2shAddresses {
		assert.Truef(t, blockchain.IsBtcP2SHAddress(address), "IsP2SHAddress should return true for address %s", address)
	}
	for _, address := range others {
		assert.Falsef(t, blockchain.IsBtcP2SHAddress(address), "IsP2SHAddress should return false for address %s", address)
	}
}

func TestIsP2WPKHAddress(t *testing.T) {
	var p2wpkhAddresses []string
	p2wpkhAddresses = append(p2wpkhAddresses, p2wpkhMainnetAddresses...)
	p2wpkhAddresses = append(p2wpkhAddresses, p2wpkhTestnetAddresses...)
	p2wpkhAddresses = append(p2wpkhAddresses, p2wpkhRegtestAddresses...)
	var others []string
	others = append(others, p2pkhMainnetAddresses...)
	others = append(others, p2pkhTestnetAddresses...)
	others = append(others, p2shMainnetAddresses...)
	others = append(others, p2shTestnetAddresses...)
	others = append(others, p2wshMainnetAddresses...)
	others = append(others, p2wshTestnetAddresses...)
	others = append(others, taprootMainnetAddresses...)
	others = append(others, taprootTestnetAddresses...)
	others = append(others, p2wshRegtestAddresses...)
	others = append(others, p2trRegtestAddresses...)

	for _, address := range p2wpkhAddresses {
		assert.Truef(t, blockchain.IsBtcP2WPKHAddress(address), "IsP2WPKHAddress should return true for address %s", address)
	}
	for _, address := range others {
		assert.Falsef(t, blockchain.IsBtcP2WPKHAddress(address), "IsP2WPKHAddress should return false for address %s", address)
	}
}

func TestIsP2WSHAddress(t *testing.T) {
	var p2wshAddresses []string
	p2wshAddresses = append(p2wshAddresses, p2wshMainnetAddresses...)
	p2wshAddresses = append(p2wshAddresses, p2wshTestnetAddresses...)
	p2wshAddresses = append(p2wshAddresses, p2wshRegtestAddresses...)
	var others []string
	others = append(others, p2pkhMainnetAddresses...)
	others = append(others, p2pkhTestnetAddresses...)
	others = append(others, p2shMainnetAddresses...)
	others = append(others, p2shTestnetAddresses...)
	others = append(others, p2wpkhMainnetAddresses...)
	others = append(others, p2wpkhTestnetAddresses...)
	others = append(others, taprootMainnetAddresses...)
	others = append(others, taprootTestnetAddresses...)
	others = append(others, p2wpkhRegtestAddresses...)
	others = append(others, p2trRegtestAddresses...)

	for _, address := range p2wshAddresses {
		assert.Truef(t, blockchain.IsBtcP2WSHAddress(address), "IsP2WSHAddress should return true for address %s", address)
	}
	for _, address := range others {
		assert.Falsef(t, blockchain.IsBtcP2WSHAddress(address), "IsP2WSHAddress should return false for address %s", address)
	}
}

func TestIsTaprootAddress(t *testing.T) {
	var taprootAddresses []string
	taprootAddresses = append(taprootAddresses, taprootMainnetAddresses...)
	taprootAddresses = append(taprootAddresses, taprootTestnetAddresses...)
	taprootAddresses = append(taprootAddresses, p2trRegtestAddresses...)
	var others []string
	others = append(others, p2pkhMainnetAddresses...)
	others = append(others, p2pkhTestnetAddresses...)
	others = append(others, p2shMainnetAddresses...)
	others = append(others, p2shTestnetAddresses...)
	others = append(others, p2wpkhMainnetAddresses...)
	others = append(others, p2wpkhTestnetAddresses...)
	others = append(others, p2wshMainnetAddresses...)
	others = append(others, p2wshTestnetAddresses...)
	others = append(others, p2wshRegtestAddresses...)
	others = append(others, p2wpkhRegtestAddresses...)

	for _, address := range taprootAddresses {
		assert.Truef(t, blockchain.IsBtcP2TRAddress(address), "IsTaprootAddress should return true for address %s", address)
	}
	for _, address := range others {
		assert.Falsef(t, blockchain.IsBtcP2TRAddress(address), "IsTaprootAddress should return false for address %s", address)
	}
}
