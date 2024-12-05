<script>
  import { toString } from "cronstrue";

  export let value = "0 0/0 * 1/1 * ? *";
  export let language = "en";
  export let showSeconds = false;
  export let showAdvanced = false;
  export let minutesStep = 1;

  let use24HourTime = true;
  let tabs = [
    // "seconds",
    "minutes",
    "hourly",
    "daily",
    "weekly",
    "monthly",
    "yearly",
    "advanced",
  ];
  export let valueTranslated = toString(value, { locale: language });
  const terminations = ["st", "nd", "rd", "th"];
  const seconds = new Array(60).fill(1).map((_el, ix) => ix);
  const minutes = new Array(60)
    .fill(1)
    .map((_el, ix) => {
      if (minutesStep > 1) {
        return ix % minutesStep === 0 ? ix : undefined;
      } else return ix;
    })
    .filter((el) => el !== undefined);
  const hours = new Array(24).fill(1).map((_el, ix) => ix);
  const days = new Array(31).fill(1).map((_el, ix) => ix + 1);
  const weekDaysShort = ["MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"];
  const weekDaysMap = {
    MON: "Monday",
    TUE: "Tuesday",
    WED: "Wednesday",
    THU: "Thursday",
    FRI: "Friday",
    SAT: "Saturday",
    SUN: "Sunday",
  };
  const weekDaysNo = [
    { key: "#1", val: "First" },
    { key: "#2", val: "Second" },
    { key: "#3", val: "Third" },
    { key: "#4", val: "Fourth" },
    { key: "#5", val: "Fifth" },
    { key: "L", val: "Last" },
  ];
  const months = new Array(12).fill(1).map((_el, ix) => ix + 1);
  const monthlyDays = [
    { key: "1W", val: "First Weekday" },
    ...days.reduce((o, el) => {
      let termination;
      const dayStr = "" + el;
      const dayStrArr = dayStr.split("");
      const lastNoOfDay = dayStrArr.pop();
      if (lastNoOfDay === "1") termination = terminations[0];
      else if (lastNoOfDay === "2") termination = terminations[1];
      else if (lastNoOfDay === "3") termination = terminations[2];
      else termination = terminations[3];
      // return el + termination + " Day";
      o.push({ key: dayStr, val: el + termination + " Day" });
      return o;
    }, []),
    { key: "LW", val: "Last Weekday" },
    { key: "L", val: "Last Day" },
  ];
  const monthNames = [
    { key: 1, val: "January" },
    { key: 2, val: "February" },
    { key: 3, val: "March" },
    { key: 4, val: "April" },
    { key: 5, val: "May" },
    { key: 6, val: "June" },
    { key: 7, val: "July" },
    { key: 8, val: "August" },
    { key: 9, val: "September" },
    { key: 10, val: "October" },
    { key: 11, val: "November" },
    { key: 12, val: "December" },
  ];

  /** @type {{
   * seconds?: { seconds?: number, },
   * minutes?: { seconds?: number, minutes?: number, },
   * hourly?: { seconds?: number, minutes?: number, hours?: number, },
   * daily?: {
   *   dailyOpt?: string,
   *   everyDay?: {
   *     seconds?: number,
   *     minutes?: number,
   *     hours?: number,
   *     hourType?: string,
   *     days?: number,
   *   },
   *   everyWeekDay?: {
   *     seconds?: number,
   *     minutes?: number,
   *     hours?: number,
   *     hourType?: string,
   *   }
   * },
   * weekly?: {
   *   MON?: boolean,
   *   TUE?: boolean,
   *   WED?: boolean,
   *   THU?: boolean,
   *   FRI?: boolean,
   *   SAT?: boolean,
   *   SUN?: boolean,
   *   seconds?: number,
   *   minutes?: number,
   *   hours?: number,
   *   hourType?: string,
   * },
   * monthly?: {
   *   monthlyOpt?: string,
   *   specificDay?: {
   *     seconds?: number,
   *     minutes?: number,
   *     hours?: number,
   *     hourType?: string,
   *     day?: string,
   *     month?: number,
   *   },
   *   specificWeekDay?: {
   *     seconds?: number,
   *     minutes?: number,
   *     hours?: number,
   *     hourType?: string,
   *     month?: number,
   *     weekDay?: string,
   *     whichWeekDay?: string,
   *   }
   * },
   * yearly?: {
   *   yearlyOpt?: string,
   *   specificMonthDay?: {
   *     seconds?: number,
   *     minutes?: number,
   *     hours?: number,
   *     hourType?: string,
   *     day?: string,
   *     month?: number,
   *   },
   *   specificMonthWeek?: {
   *     seconds?: number,
   *     minutes?: number,
   *     hours?: number,
   *     hourType?: string,
   *     month?: number,
   *     weekDay?: string,
   *     whichWeekDay?: string,
   *   }
   * },
   * advanced?: string,
   * }} cronSetts
   */
  let cronSetts = {
    seconds: {
      seconds: 0,
    },
    minutes: {
      seconds: 0,
      minutes: 0,
    },
    hourly: {
      seconds: 0,
      minutes: 0,
      hours: 0,
    },
    daily: {
      dailyOpt: "everyDay",
      everyDay: {
        seconds: 0,
        minutes: 0,
        hours: 0,
        hourType: "AM",
        days: 1,
      },
      everyWeekDay: {
        seconds: 0,
        minutes: 0,
        hours: 0,
        hourType: "AM",
      },
    },
    weekly: {
      MON: true,
      TUE: false,
      WED: false,
      THU: false,
      FRI: false,
      SAT: false,
      SUN: false,
      seconds: 0,
      minutes: 0,
      hours: 0,
      hourType: "AM",
    },
    monthly: {
      monthlyOpt: "specificDay",
      specificDay: {
        seconds: 0,
        minutes: 0,
        hours: 0,
        hourType: "AM",
        day: "1W",
        month: 1,
      },
      specificWeekDay: {
        seconds: 0,
        minutes: 0,
        hours: 0,
        hourType: "AM",
        month: 1,
        weekDay: "MON",
        whichWeekDay: "#1",
      },
    },
    yearly: {
      yearlyOpt: "specificMonthDay",
      specificMonthDay: {
        seconds: 0,
        minutes: 0,
        hours: 0,
        hourType: "AM",
        day: "1W",
        month: 1,
      },
      specificMonthWeek: {
        seconds: 0,
        minutes: 0,
        hours: 0,
        hourType: "AM",
        month: 1,
        weekDay: "MON",
        whichWeekDay: "#1",
      },
    },
    advanced: null,
  };
  let cronSettsBkup = JSON.stringify(cronSetts);
  let tabShow = "minutes";
  let valueBkUp = value;

  parseCron();

  function generateCron() {
    if ("minutes" === tabShow) {
      value = `${cronSetts.minutes.seconds} 0/${cronSetts.minutes.minutes} * 1/1 * ? *`;
    }
    if ("hourly" === tabShow) {
      value = `${cronSetts.hourly.seconds} ${cronSetts.hourly.minutes} 0/${cronSetts.hourly.hours} 1/1 * ? *`;
    }
    if ("daily" === tabShow) {
      if (cronSetts.daily.dailyOpt === "everyDay") {
        value = `${cronSetts.daily.everyDay.seconds} ${
          cronSetts.daily.everyDay.minutes
        } ${hourToCron(
          cronSetts.daily.everyDay.hours,
          cronSetts.daily.everyDay.hourType
        )} 1/${cronSetts.daily.everyDay.days} * ? *`;
      }
      if (cronSetts.daily.dailyOpt === "everyWeekDay") {
        value = `${cronSetts.daily.everyWeekDay.seconds} ${
          cronSetts.daily.everyWeekDay.minutes
        } ${hourToCron(
          cronSetts.daily.everyWeekDay.hours,
          cronSetts.daily.everyWeekDay.hourType
        )} ? * MON-FRI *`;
      }
    }
    if ("weekly" === tabShow) {
      const days = Object.keys(cronSetts.weekly)
        .reduce((o, key) => {
          if (weekDaysShort.includes(key) && cronSetts.weekly[key]) {
            o = [...o, key];
          }
          return o;
        }, [])
        .join(",");
      value = `${cronSetts.weekly.seconds} ${
        cronSetts.weekly.minutes
      } ${hourToCron(
        cronSetts.weekly.hours,
        cronSetts.weekly.hourType
      )} ? * ${days} *`;
    }
    if ("monthly" === tabShow) {
      if ("specificDay" === cronSetts.monthly.monthlyOpt) {
        value = `${cronSetts.monthly.specificDay.seconds} ${
          cronSetts.monthly.specificDay.minutes
        } ${hourToCron(
          cronSetts.monthly.specificDay.hours,
          cronSetts.monthly.specificDay.hourType
        )} ${cronSetts.monthly.specificDay.day} 1/${
          cronSetts.monthly.specificDay.month
        } ? *`;
      }
      if ("specificWeekDay" === cronSetts.monthly.monthlyOpt) {
        value = `${cronSetts.monthly.specificWeekDay.seconds} ${
          cronSetts.monthly.specificWeekDay.minutes
        } ${hourToCron(
          cronSetts.monthly.specificWeekDay.hours,
          cronSetts.monthly.specificWeekDay.hourType
        )} ? 1/${cronSetts.monthly.specificWeekDay.month} ${
          cronSetts.monthly.specificWeekDay.weekDay
        }${cronSetts.monthly.specificWeekDay.whichWeekDay} *`;
      }
    }
    if ("yearly" === tabShow) {
      if ("specificMonthDay" === cronSetts.yearly.yearlyOpt) {
        value = `${cronSetts.yearly.specificMonthDay.seconds} ${
          cronSetts.yearly.specificMonthDay.minutes
        } ${hourToCron(
          cronSetts.yearly.specificMonthDay.hours,
          cronSetts.yearly.specificMonthDay.hourType
        )} ${cronSetts.yearly.specificMonthDay.day} ${
          cronSetts.yearly.specificMonthDay.month
        } ? *`;
      }
      if ("specificMonthWeek" === cronSetts.yearly.yearlyOpt) {
        value = `${cronSetts.yearly.specificMonthWeek.seconds} ${
          cronSetts.yearly.specificMonthWeek.minutes
        } ${hourToCron(
          cronSetts.yearly.specificMonthWeek.hours,
          cronSetts.yearly.specificMonthWeek.hourType
        )} ? ${cronSetts.yearly.specificMonthWeek.month} ${
          cronSetts.yearly.specificMonthWeek.weekDay
        }${cronSetts.yearly.specificMonthWeek.whichWeekDay} *`;
      }
    }
    if ("advanced" === tabShow) {
      value = cronSetts.advanced;
    }
    translateCron();
  }

  function translateCron() {
    valueTranslated = toString(value, { locale: language });
  }

  /**
   * @param {string} which
   */
  function changeTab(which) {
    tabShow = which;
    value = valueBkUp;
    cronSetts = JSON.parse(cronSettsBkup);
    generateCron();
    translateCron();
  }

  /**
   * @param {number} hours
   */
  function processHour(hours) {
    if (use24HourTime) {
      return hours;
    } else {
      return ((hours + 11) % 12) + 1;
    }
  }

  /**
   * @param {number} hours
   */
  function getHourType(hours) {
    return use24HourTime ? null : hours >= 12 ? "PM" : "AM";
  }

  function parseCron() {
    const cron = value;
    const segments = cron.split(" ");
    if (segments.length === 6 || segments.length === 7) {
      const [seconds, minutes, hours, dayOfMonth, month, dayOfWeek] = segments;
      if (cron.match(/\d+ 0\/\d+ \* 1\/1 \* \? \*/)) {
        tabShow = "minutes";
        cronSetts.minutes.minutes = parseInt(minutes.substring(2));
        cronSetts.minutes.seconds = parseInt(seconds);
      } else if (cron.match(/\d+ \d+ 0\/\d+ 1\/1 \* \? \*/)) {
        tabShow = "hourly";
        cronSetts.hourly.hours = parseInt(hours.substring(2));
        cronSetts.hourly.minutes = parseInt(minutes);
        cronSetts.hourly.seconds = parseInt(seconds);
      } else if (cron.match(/\d+ \d+ \d+ 1\/\d+ \* \? \*/)) {
        tabShow = "daily";
        cronSetts.daily.dailyOpt = "everyDay";
        cronSetts.daily.everyDay.days = parseInt(dayOfMonth.substring(2));
        const parsedHours = parseInt(hours);
        cronSetts.daily.everyDay.hours = processHour(parsedHours);
        cronSetts.daily.everyDay.hourType = getHourType(parsedHours);
        cronSetts.daily.everyDay.minutes = parseInt(minutes);
        cronSetts.daily.everyDay.seconds = parseInt(seconds);
      } else if (cron.match(/\d+ \d+ \d+ \? \* MON-FRI \*/)) {
        tabShow = "daily";
        cronSetts.daily.dailyOpt = "everyWeekDay";
        const parsedHours = parseInt(hours);
        cronSetts.daily.everyWeekDay.hours = processHour(parsedHours);
        cronSetts.daily.everyWeekDay.hourType = getHourType(parsedHours);
        cronSetts.daily.everyWeekDay.minutes = parseInt(minutes);
        cronSetts.daily.everyWeekDay.seconds = parseInt(seconds);
      } else if (
        cron.match(
          /\d+ \d+ \d+ \? \* (MON|TUE|WED|THU|FRI|SAT|SUN)(,(MON|TUE|WED|THU|FRI|SAT|SUN))* \*/
        )
      ) {
        tabShow = "weekly";
        Object.keys(cronSetts.weekly).forEach((key) => {
          if (weekDaysShort.includes(key)) cronSetts.weekly[key] = false;
        });
        dayOfWeek
          .split(",")
          .forEach((weekDay) => (cronSetts.weekly[weekDay] = true));
        const parsedHours = parseInt(hours);
        cronSetts.weekly.hours = processHour(parsedHours);
        cronSetts.weekly.hourType = getHourType(parsedHours);
        cronSetts.weekly.minutes = parseInt(minutes);
        cronSetts.weekly.seconds = parseInt(seconds);
      } else if (cron.match(/\d+ \d+ \d+ (\d+|L|LW|1W) 1\/\d+ \? \*/)) {
        tabShow = "monthly";
        cronSetts.monthly.monthlyOpt = "specificDay";
        cronSetts.monthly.specificDay.day = dayOfMonth;
        cronSetts.monthly.specificDay.month = parseInt(month.substring(2));
        const parsedHours = parseInt(hours);
        cronSetts.monthly.specificDay.hours = processHour(parsedHours);
        cronSetts.monthly.specificDay.hourType = getHourType(parsedHours);
        cronSetts.monthly.specificDay.minutes = parseInt(minutes);
        cronSetts.monthly.specificDay.seconds = parseInt(seconds);
      } else if (
        cron.match(
          /\d+ \d+ \d+ \? 1\/\d+ (MON|TUE|WED|THU|FRI|SAT|SUN)((#[1-5])|L) \*/
        )
      ) {
        const day = dayOfWeek.substr(0, 3);
        const monthWeek = dayOfWeek.substr(3);
        tabShow = "monthly";
        cronSetts.monthly.monthlyOpt = "specificWeekDay";
        cronSetts.monthly.specificWeekDay.whichWeekDay = monthWeek;
        cronSetts.monthly.specificWeekDay.weekDay = day;
        cronSetts.monthly.specificWeekDay.month = parseInt(month.substring(2));
        const parsedHours = parseInt(hours);
        cronSetts.monthly.specificWeekDay.hours = processHour(parsedHours);
        cronSetts.monthly.specificWeekDay.hourType = getHourType(parsedHours);
        cronSetts.monthly.specificWeekDay.minutes = parseInt(minutes);
        cronSetts.monthly.specificWeekDay.seconds = parseInt(seconds);
      } else if (cron.match(/\d+ \d+ \d+ (\d+|L|LW|1W) \d+ \? \*/)) {
        tabShow = "yearly";
        cronSetts.yearly.yearlyOpt = "specificMonthDay";
        cronSetts.yearly.specificMonthDay.month = parseInt(month);
        cronSetts.yearly.specificMonthDay.day = dayOfMonth;
        const parsedHours = parseInt(hours);
        cronSetts.yearly.specificMonthDay.hours = processHour(parsedHours);
        cronSetts.yearly.specificMonthDay.hourType = getHourType(parsedHours);
        cronSetts.yearly.specificMonthDay.minutes = parseInt(minutes);
        cronSetts.yearly.specificMonthDay.seconds = parseInt(seconds);
      } else if (
        cron.match(
          /\d+ \d+ \d+ \? \d+ (MON|TUE|WED|THU|FRI|SAT|SUN)((#[1-5])|L) \*/
        )
      ) {
        const day = dayOfWeek.substr(0, 3);
        const monthWeek = dayOfWeek.substr(3);
        tabShow = "yearly";
        cronSetts.yearly.yearlyOpt = "specificMonthWeek";
        cronSetts.yearly.specificMonthWeek.whichWeekDay = monthWeek;
        cronSetts.yearly.specificMonthWeek.weekDay = day;
        cronSetts.yearly.specificMonthWeek.month = parseInt(month);
        const parsedHours = parseInt(hours);
        cronSetts.yearly.specificMonthWeek.hours = processHour(parsedHours);
        cronSetts.yearly.specificMonthWeek.hourType = getHourType(parsedHours);
        cronSetts.yearly.specificMonthWeek.minutes = parseInt(minutes);
        cronSetts.yearly.specificMonthWeek.seconds = parseInt(seconds);
      } else {
        if (showAdvanced) {
          tabShow = "advanced";
          cronSetts.advanced = cron;
        } else {
          tabShow = "minutes";
        }
      }
    } else {
      throw new Error(
        "Unsupported cron expression. Expression must be 6 or 7 segments"
      );
    }
  }

  /**
   * @param {number} hour
   * @param {string} hourType
   */
  function hourToCron(hour, hourType) {
    if (use24HourTime) {
      return hour;
    } else {
      return hourType === "AM"
        ? hour === 12
          ? 0
          : hour
        : hour === 12
        ? 12
        : hour + 12;
    }
  }
</script>

<div>
	<h2>{valueTranslated}</h2>
</div>

<div>
  {#each tabs as tab}
    {#if showAdvanced || (!showAdvanced && tab !== "advanced")}
      <button
        class={tab === tabShow ? "active" : ""}
        on:click={() => changeTab(tab)}>{tab}</button
      >
    {/if}
  {/each}
</div>
<hr />

{#if showAdvanced && tabShow === "advanced"}
  <div class="advanced">
    <input
      type="text"
      bind:value={cronSetts.advanced}
      on:change={generateCron}
    />
  </div>
{/if}

{#if tabShow === "yearly"}
  <div class="yearly">
    <div>
      <input
        type="radio"
        name="yearly-opts"
        bind:group={cronSetts.yearly.yearlyOpt}
        value="specificMonthDay"
        on:change={() => generateCron()}
      />
      Every
      <select
        bind:value={cronSetts.yearly.specificMonthDay.month}
        on:change={() => generateCron()}
      >
        {#each monthNames as item}
          <option value={item.key}>{item.val}</option>
        {/each}
      </select>
      on the
      <select
        bind:value={cronSetts.yearly.specificMonthDay.day}
        on:change={() => generateCron()}
      >
        {#each monthlyDays as item}
          <option value={item.key}>{item.val}</option>
        {/each}
      </select>
      at
      <select
        bind:value={cronSetts.yearly.specificMonthDay.hours}
        on:change={() => generateCron()}
      >
        {#each hours as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      <select
        bind:value={cronSetts.yearly.specificMonthDay.minutes}
        on:change={() => generateCron()}
      >
        {#each minutes as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      {#if showSeconds}
        <select
          bind:value={cronSetts.yearly.specificMonthDay.seconds}
          on:change={() => generateCron()}
        >
          {#each seconds as item}
            <option value={item}>{item}</option>
          {/each}
        </select>
      {/if}
    </div>
    <div>
      <input
        type="radio"
        name="yearly-opts"
        bind:group={cronSetts.yearly.yearlyOpt}
        value="specificMonthWeek"
        on:change={() => generateCron()}
      />
      On the
      <select
        bind:value={cronSetts.yearly.specificMonthWeek.whichWeekDay}
        on:change={() => generateCron()}
      >
        {#each weekDaysNo as item}
          <option value={item.key}>{item.val}</option>
        {/each}
      </select>
      <select
        bind:value={cronSetts.yearly.specificMonthWeek.weekDay}
        on:change={() => generateCron()}
      >
        {#each Object.entries(weekDaysMap) as [key, val]}
          <option value={key}>{val}</option>
        {/each}
      </select>
      of
      <select
        bind:value={cronSetts.yearly.specificMonthWeek.month}
        on:change={() => generateCron()}
      >
        {#each monthNames as item}
          <option value={item.key}>{item.val}</option>
        {/each}
      </select>
      at
      <select
        bind:value={cronSetts.yearly.specificMonthWeek.hours}
        on:change={() => generateCron()}
      >
        {#each hours as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      <select
        bind:value={cronSetts.yearly.specificMonthWeek.minutes}
        on:change={() => generateCron()}
      >
        {#each minutes as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      {#if showSeconds}
        <select
          bind:value={cronSetts.yearly.specificMonthWeek.seconds}
          on:change={() => generateCron()}
        >
          {#each seconds as item}
            <option value={item}>{item}</option>
          {/each}
        </select>
      {/if}
    </div>
  </div>
{/if}

{#if tabShow === "monthly"}
  <div class="monthly">
    <div>
      <input
        type="radio"
        name="monthly-opts"
        bind:group={cronSetts.monthly.monthlyOpt}
        value="specificDay"
        on:change={() => generateCron()}
      />
      On the
      <select
        bind:value={cronSetts.monthly.specificDay.day}
        on:change={() => generateCron()}
      >
        {#each monthlyDays as item}
          <option value={item.key}>{item.val}</option>
        {/each}
      </select>
      of every
      <select
        bind:value={cronSetts.monthly.specificDay.month}
        on:change={() => generateCron()}
      >
        {#each months as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      month(s) at
      <select
        bind:value={cronSetts.monthly.specificDay.hours}
        on:change={() => generateCron()}
      >
        {#each hours as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      <select
        bind:value={cronSetts.monthly.specificDay.minutes}
        on:change={() => generateCron()}
      >
        {#each minutes as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      {#if showSeconds}
        <select
          bind:value={cronSetts.monthly.specificDay.seconds}
          on:change={() => generateCron()}
        >
          {#each seconds as item}
            <option value={item}>{item}</option>
          {/each}
        </select>
      {/if}
    </div>
    <div>
      <input
        type="radio"
        name="monthly-opts"
        bind:group={cronSetts.monthly.monthlyOpt}
        value="specificWeekDay"
        on:change={() => generateCron()}
      />
      On the
      <select
        bind:value={cronSetts.monthly.specificWeekDay.whichWeekDay}
        on:change={() => generateCron()}
      >
        {#each weekDaysNo as item}
          <option value={item.key}>{item.val}</option>
        {/each}
      </select>
      <select
        bind:value={cronSetts.monthly.specificWeekDay.weekDay}
        on:change={() => generateCron()}
      >
        {#each weekDaysShort as item}
          <option value={item}>{weekDaysMap[item]}</option>
        {/each}
      </select>
      of every
      <select
        bind:value={cronSetts.monthly.specificWeekDay.month}
        on:change={() => generateCron()}
      >
        {#each months as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      at
      <select
        bind:value={cronSetts.monthly.specificWeekDay.hours}
        on:change={() => generateCron()}
      >
        {#each hours as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      <select
        bind:value={cronSetts.monthly.specificWeekDay.minutes}
        on:change={() => generateCron()}
      >
        {#each minutes as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      {#if showSeconds}
        <select
          bind:value={cronSetts.monthly.specificWeekDay.seconds}
          on:change={() => generateCron()}
        >
          {#each seconds as item}
            <option value={item}>{item}</option>
          {/each}
        </select>
      {/if}
    </div>
  </div>
{/if}

{#if tabShow === "weekly"}
  <div class="weekly">
    {#each weekDaysShort as item}
      <label>
        <input
          type="checkbox"
          value={item}
          bind:checked={cronSetts.weekly[item]}
          on:change={() => generateCron()}
        />
        {weekDaysMap[item]}
      </label>
    {/each}
    <div>
      Start time
      <select
        bind:value={cronSetts.weekly.hours}
        on:change={() => generateCron()}
      >
        {#each hours as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      <select
        bind:value={cronSetts.weekly.minutes}
        on:change={() => generateCron()}
      >
        {#each minutes as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      {#if showSeconds}
        <select
          bind:value={cronSetts.weekly.seconds}
          on:change={() => generateCron()}
        >
          {#each seconds as item}
            <option value={item}>{item}</option>
          {/each}
        </select>
      {/if}
    </div>
  </div>
{/if}

{#if tabShow === "daily"}
  <div class="daily">
    <div class="entity">
      <input
        type="radio"
        name="daily-opt"
        bind:group={cronSetts.daily.dailyOpt}
        value="everyDay"
        on:change={() => generateCron()}
      />
      Every
      <select
        bind:value={cronSetts.daily.everyDay.days}
        on:change={() => generateCron()}
      >
        {#each days as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      day(s)
    </div>
    <div class="entity">
      at
      <select
        bind:value={cronSetts.daily.everyDay.hours}
        on:change={() => generateCron()}
      >
        {#each hours as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      <select
        bind:value={cronSetts.daily.everyDay.minutes}
        on:change={() => generateCron()}
      >
        {#each minutes as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      {#if showSeconds}
        <select
          bind:value={cronSetts.daily.everyDay.seconds}
          on:change={() => generateCron()}
        >
          {#each seconds as item}
            <option value={item}>{item}</option>
          {/each}
        </select>
      {/if}
    </div>

    <div>
      <input
        type="radio"
        name="daily-opt"
        bind:group={cronSetts.daily.dailyOpt}
        value="everyWeekDay"
        on:change={() => generateCron()}
      />
      Every week day (Monday through Friday) at
      <select
        bind:value={cronSetts.daily.everyWeekDay.hours}
        on:change={() => generateCron()}
      >
        {#each hours as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      <select
        bind:value={cronSetts.daily.everyWeekDay.minutes}
        on:change={() => generateCron()}
      >
        {#each minutes as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      {#if showSeconds}
        <select
          bind:value={cronSetts.daily.everyWeekDay.seconds}
          on:change={() => generateCron()}
        >
          {#each seconds as item}
            <option value={item}>{item}</option>
          {/each}
        </select>
      {/if}
    </div>
  </div>
{/if}

{#if tabShow === "hourly"}
  <div class="hourly">
    <div class="entity">
      Every
      <select
        bind:value={cronSetts.hourly.hours}
        on:change={() => generateCron()}
      >
        {#each hours as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      hour(s)
    </div>
    <div class="entity">
      on minute
      <select
        bind:value={cronSetts.hourly.minutes}
        on:change={() => generateCron()}
      >
        {#each minutes as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
    </div>
    {#if showSeconds}
      <div class="entity">
        and second
        <select
          bind:value={cronSetts.hourly.seconds}
          on:change={() => generateCron()}
        >
          {#each seconds as item}
            <option value={item}>{item}</option>
          {/each}
        </select>
      </div>
    {/if}
  </div>
{/if}

{#if tabShow === "minutes"}
  <div class="minutes">
    <div class="entity">
      Every
      <select
        bind:value={cronSetts.minutes.minutes}
        on:change={() => generateCron()}
      >
        {#each minutes as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      minute(s)
    </div>
    {#if showSeconds}
      <div class="entity">
        on second
        <select
          bind:value={cronSetts.minutes.seconds}
          on:change={() => generateCron()}
        >
          {#each seconds as item}
            <option value={item}>{item}</option>
          {/each}
        </select>
      </div>
    {/if}
  </div>
{/if}

{#if tabShow === "seconds"}
  <div class="seconds">
    <div>
      Every
      <select
        bind:value={cronSetts.seconds.seconds}
        on:change={() => generateCron()}
      >
        {#each seconds as item}
          <option value={item}>{item}</option>
        {/each}
      </select>
      second(s)
    </div>
  </div>
{/if}

<hr />

<!-- <hr /> -->
<!-- <pre> -->
<!-- {JSON.stringify(cronSetts, null, 2)} -->

<!-- </pre> -->
<style>
	hr {
		color: var(--trough-color);
	}
	button {
		margin: 2px;
		padding: 8px;
		border: none;
		border-radius: 6px;
	}
	.active {
		background-color: var(--accent-bg);
		color: var(--bg);
	}
	.entity {
		display: inline-block;
	}
	.weekly {
		gap: 10px;
	}
</style>
